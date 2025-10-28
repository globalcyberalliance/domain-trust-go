package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/klauspost/compress/zstd"
	"github.com/moul/http2curl"
	"github.com/spf13/cast"
)

const (
	ContentTypeCBOR  = "application/cbor"
	ContentTypeJSON  = "application/json"
	EncodingTypeGZIP = "gzip"
	EncodingTypeZSTD = "zstd"
)

type (
	GenericResponse struct {
		Status int    `json:"status"`
		Detail string `json:"detail"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
)

func (r GenericResponse) ToErrorString() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s (HTTP %d)", r.Detail, r.Status))

	if len(r.Errors) > 0 {
		sb.WriteString(": ")
		for i, e := range r.Errors {
			sb.WriteString(e.Message)
			if i < len(r.Errors)-1 {
				sb.WriteString("; ")
			}
		}
	}

	return sb.String()
}

func (c *Client) DELETE(ctx context.Context, endpoint string, obj any) ([]byte, error) {
	return c.makeRequest(ctx, endpoint, "DELETE", nil, obj)
}

func (c *Client) GET(ctx context.Context, endpoint string, obj any) ([]byte, error) {
	return c.makeRequest(ctx, endpoint, "GET", nil, obj)
}

func (c *Client) PATCH(ctx context.Context, endpoint string, body []byte, obj any) ([]byte, error) {
	return c.makeRequest(ctx, endpoint, "PATCH", body, obj)
}

func (c *Client) POST(ctx context.Context, endpoint string, body []byte, obj any) ([]byte, error) {
	return c.makeRequest(ctx, endpoint, "POST", body, obj)
}

func (c *Client) makeRequest(ctx context.Context, endpoint string, method string, requestBody []byte, object any) ([]byte, error) {
	endpointURL := fmt.Sprintf("%s/%s", EndpointURL, endpoint)

	if len(requestBody) > 0 {
		var compressedBody bytes.Buffer

		encoder, err := zstd.NewWriter(&compressedBody)
		if err != nil {
			return nil, fmt.Errorf("create compression writer: %w", err)
		}
		defer encoder.Close()

		if _, err = encoder.Write(requestBody); err != nil {
			return nil, fmt.Errorf("write compressed data: %w", err)
		}

		if err = encoder.Close(); err != nil {
			return nil, fmt.Errorf("close compression encoder: %w", err)
		}

		requestBody = compressedBody.Bytes()
	}

	req, err := http.NewRequestWithContext(ctx, method, endpointURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Add("Authorization", "Bearer "+c.apiKey)
	}

	req.Header.Add("Accept", c.contentType)
	req.Header.Add("Accept-Encoding", c.encodingType)
	req.Header.Add("Content-Encoding", c.encodingType)
	req.Header.Add("Content-Type", c.contentType)
	req.Header.Add("DT-Client-Version", Version)

	if c.debug {
		curl, cErr := http2curl.GetCurlCommand(req)
		if cErr == nil {
			fmt.Println(curl.String())
		}
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	if res == nil {
		return nil, errors.New("read response")
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	resBody, err = unmarshalResBody(resBody, res.Header.Get("Content-Encoding"))
	if err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	if res.StatusCode >= http.StatusBadRequest {
		if c.debug {
			fmt.Println(string(resBody))
		}

		if len(resBody) > 0 {
			resp := GenericResponse{}

			switch res.Header.Get("Content-Type") {
			case ContentTypeCBOR:
				err = cbor.Unmarshal(resBody, &resp)
			case ContentTypeJSON:
				err = json.Unmarshal(resBody, &resp)
			}

			if err == nil {
				return nil, errors.New(resp.ToErrorString())
			}
		}

		// Return the body as it may contain a useful error message.
		return resBody, errors.New("request status code " + cast.ToString(res.StatusCode))
	}

	if len(resBody) > 0 && object != nil {
		switch res.Header.Get("Content-Type") {
		case ContentTypeCBOR:
			err = cbor.Unmarshal(resBody, object)
		case ContentTypeJSON:
			err = json.Unmarshal(resBody, object)
		}
		if err != nil {
			return nil, fmt.Errorf("request succeeded, couldn't unmarshal into object: %w", err)
		}
	}

	return resBody, nil
}

// parseTag splits a struct tag like `foo:"name,omitempty,other"` into
// ("name", map["omitempty"]=true, map["other"]=true).
func parseTag(tag string) (string, map[string]bool) {
	name := ""
	opts := map[string]bool{}
	if tag == "" {
		return name, opts
	}

	parts := strings.Split(tag, ",")
	name = parts[0]
	for _, p := range parts[1:] {
		if p != "" {
			opts[p] = true
		}
	}

	return name, opts
}

func structToQueryParams(data interface{}) string {
	values := url.Values{}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}

	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i)

		// Key and value parts in a query string pair.
		const keyValuePairParts = 2

		// Recurse into anonymous structs to surface their fields.
		if field.Anonymous && value.Kind() == reflect.Struct {
			// Flatten embedded struct.
			embeddedQS := structToQueryParams(value.Interface())
			if embeddedQS != "" {
				for _, kv := range strings.Split(embeddedQS, "&") {
					if kv == "" {
						continue
					}

					parts := strings.SplitN(kv, "=", keyValuePairParts)
					if len(parts) == keyValuePairParts {
						values.Add(parts[0], parts[1])
					}
				}
			}
			continue
		}

		// Read the `query` tag (primary source of the parameter name)
		qName, _ := parseTag(field.Tag.Get("query"))
		if qName == "" || qName == "-" {
			continue
		}

		// Also read the `json` options to honor omitempty there too (optional).
		_, jsonOpts := parseTag(field.Tag.Get("json"))
		omitEmpty := jsonOpts["omitempty"] || jsonOpts["omitzero"]

		// Zero checks.
		switch val := value.Interface().(type) {
		case netip.Addr:
			if omitEmpty && !val.IsValid() {
				continue
			}

			values.Set(qName, val.String())
		case time.Time:
			if omitEmpty && val.IsZero() {
				continue
			}

			values.Set(qName, val.Format(time.RFC3339))
		default:
			// Use IsZero when available; handles strings "", numbers 0, time.Time{}, nil pointers/slices/maps, etc.
			if omitEmpty && value.IsZero() {
				continue
			}

			values.Set(qName, fmt.Sprintf("%v", val))
		}
	}

	return values.Encode()
}

func unmarshalResBody(body []byte, encodingType string) ([]byte, error) {
	if len(body) == 0 || encodingType == "" {
		return body, nil
	}

	var resBody []byte

	switch encodingType {
	case EncodingTypeGZIP:
		gzipReader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("decode gzip content: %w", err)
		}
		defer gzipReader.Close()

		resBody, err = io.ReadAll(gzipReader)
		if err != nil {
			return nil, fmt.Errorf("read decompressed gzip body: %w", err)
		}
	case EncodingTypeZSTD:
		zstdReader, err := zstd.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("decode zstd contentType: %w", err)
		}
		defer zstdReader.Close()

		resBody, err = io.ReadAll(zstdReader)
		if err != nil {
			return nil, fmt.Errorf("read decompressed body: %w", err)
		}
	}

	return resBody, nil
}
