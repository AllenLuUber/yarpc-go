package json

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"go.uber.org/yarpc"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/api/transport/transporttest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type simpleRequest struct {
	Name       string
	Attributes map[string]int32
}

type simpleResponse struct {
	Success bool
}

func TestHandleStructSuccess(t *testing.T) {
	h := func(ctx context.Context, r yarpc.ReqMeta, body *simpleRequest) (*simpleResponse, yarpc.ResMeta, error) {
		assert.Equal(t, "simpleCall", r.Procedure())
		assert.Equal(t, "foo", body.Name)
		assert.Equal(t, map[string]int32{"bar": 42}, body.Attributes)

		return &simpleResponse{Success: true}, nil, nil
	}

	handler := jsonHandler{
		reader:  structReader{reflect.TypeOf(simpleRequest{})},
		handler: reflect.ValueOf(h),
	}

	resw := new(transporttest.FakeResponseWriter)
	err := handler.Handle(context.Background(), &transport.Request{
		Procedure: "simpleCall",
		Encoding:  "json",
		Body:      jsonBody(`{"name": "foo", "attributes": {"bar": 42}}`),
	}, resw)
	require.NoError(t, err)

	var response simpleResponse
	require.NoError(t, json.Unmarshal(resw.Body.Bytes(), &response))

	assert.Equal(t, simpleResponse{Success: true}, response)
}

func TestHandleMapSuccess(t *testing.T) {
	h := func(ctx context.Context, _ yarpc.ReqMeta, body map[string]interface{}) (map[string]string, yarpc.ResMeta, error) {
		assert.Equal(t, 42.0, body["foo"])
		assert.Equal(t, []interface{}{"a", "b", "c"}, body["bar"])

		return map[string]string{"success": "true"}, nil, nil
	}

	handler := jsonHandler{
		reader:  mapReader{reflect.TypeOf(make(map[string]interface{}))},
		handler: reflect.ValueOf(h),
	}

	resw := new(transporttest.FakeResponseWriter)
	err := handler.Handle(context.Background(), &transport.Request{
		Procedure: "foo",
		Encoding:  "json",
		Body:      jsonBody(`{"foo": 42, "bar": ["a", "b", "c"]}`),
	}, resw)
	require.NoError(t, err)

	var response struct{ Success string }
	require.NoError(t, json.Unmarshal(resw.Body.Bytes(), &response))
	assert.Equal(t, "true", response.Success)
}

func TestHandleInterfaceEmptySuccess(t *testing.T) {
	h := func(ctx context.Context, _ yarpc.ReqMeta, body interface{}) (interface{}, yarpc.ResMeta, error) {
		return body, nil, nil
	}

	handler := jsonHandler{reader: ifaceEmptyReader{}, handler: reflect.ValueOf(h)}

	resw := new(transporttest.FakeResponseWriter)
	err := handler.Handle(context.Background(), &transport.Request{
		Procedure: "foo",
		Encoding:  "json",
		Body:      jsonBody(`["a", "b", "c"]`),
	}, resw)
	require.NoError(t, err)

	assert.JSONEq(t, `["a", "b", "c"]`, resw.Body.String())
}

func TestHandleSuccessWithResponseHeaders(t *testing.T) {
	h := func(ctx context.Context, r yarpc.ReqMeta, _ *simpleRequest) (*simpleResponse, yarpc.ResMeta, error) {
		resMeta := yarpc.NewResMeta().Headers(yarpc.NewHeaders().With("foo", "bar"))
		return &simpleResponse{Success: true}, resMeta, nil
	}

	handler := jsonHandler{
		reader:  structReader{reflect.TypeOf(simpleRequest{})},
		handler: reflect.ValueOf(h),
	}

	resw := new(transporttest.FakeResponseWriter)
	err := handler.Handle(context.Background(), &transport.Request{
		Procedure: "simpleCall",
		Encoding:  "json",
		Body:      jsonBody(`{"name": "foo", "attributes": {"bar": 42}}`),
	}, resw)
	require.NoError(t, err)

	assert.Equal(t, transport.NewHeaders().With("foo", "bar"), resw.Headers)
}

func jsonBody(s string) io.Reader {
	return bytes.NewReader([]byte(s))
}
