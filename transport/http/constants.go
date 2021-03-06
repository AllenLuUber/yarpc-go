// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package http

const (
	// ApplicationHeaderPrefix is the prefix added to application headers over
	// the wire.
	ApplicationHeaderPrefix = "Rpc-Header-"

	// TODO(abg): Allow customizing header prefixes

	// CallerHeader is the HTTP header used to indiate the service doing the calling
	CallerHeader = "Rpc-Caller"

	// EncodingHeader is the HTTP header used to specify the name of the
	// encoding.
	EncodingHeader = "Rpc-Encoding"

	// TTLMSHeader is the HTTP header used to indicate the ttl in ms
	TTLMSHeader = "Context-TTL-MS"

	// ProcedureHeader is the HTTP header used to indicate the procedure
	ProcedureHeader = "Rpc-Procedure"

	// ServiceHeader is the HTTP header used to indicate the service
	ServiceHeader = "Rpc-Service"

	// ShardKeyHeader is the HTTP header used by a clustered service to
	// identify the node responsible for handling the request.
	ShardKeyHeader = "Rpc-Shard-Key"

	// RoutingKeyHeader is the HTTP header used by a relay to identify
	// the traffic group responsible for handling the request, overriding the
	// service name when available.
	RoutingKeyHeader = "Rpc-Routing-Key"

	// RoutingDelegateHeader is the HTTP header used by a relay to identify a
	// service that can proxy for the destined service and perform
	// application-specific routing.
	RoutingDelegateHeader = "Rpc-Routing-Delegate"
)
