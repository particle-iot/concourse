// Code generated by counterfeiter. DO NOT EDIT.
package connectionfakes

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/url"
	"sync"

	"github.com/concourse/concourse/atc/worker/gardenruntime/gclient/connection"
	"github.com/tedsuo/rata"
)

type FakeHijackStreamer struct {
	HijackStub        func(context.Context, string, io.Reader, rata.Params, url.Values, string) (net.Conn, *bufio.Reader, error)
	hijackMutex       sync.RWMutex
	hijackArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 io.Reader
		arg4 rata.Params
		arg5 url.Values
		arg6 string
	}
	hijackReturns struct {
		result1 net.Conn
		result2 *bufio.Reader
		result3 error
	}
	hijackReturnsOnCall map[int]struct {
		result1 net.Conn
		result2 *bufio.Reader
		result3 error
	}
	StreamStub        func(string, io.Reader, rata.Params, url.Values, string) (io.ReadCloser, error)
	streamMutex       sync.RWMutex
	streamArgsForCall []struct {
		arg1 string
		arg2 io.Reader
		arg3 rata.Params
		arg4 url.Values
		arg5 string
	}
	streamReturns struct {
		result1 io.ReadCloser
		result2 error
	}
	streamReturnsOnCall map[int]struct {
		result1 io.ReadCloser
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeHijackStreamer) Hijack(arg1 context.Context, arg2 string, arg3 io.Reader, arg4 rata.Params, arg5 url.Values, arg6 string) (net.Conn, *bufio.Reader, error) {
	fake.hijackMutex.Lock()
	ret, specificReturn := fake.hijackReturnsOnCall[len(fake.hijackArgsForCall)]
	fake.hijackArgsForCall = append(fake.hijackArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 io.Reader
		arg4 rata.Params
		arg5 url.Values
		arg6 string
	}{arg1, arg2, arg3, arg4, arg5, arg6})
	stub := fake.HijackStub
	fakeReturns := fake.hijackReturns
	fake.recordInvocation("Hijack", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.hijackMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeHijackStreamer) HijackCallCount() int {
	fake.hijackMutex.RLock()
	defer fake.hijackMutex.RUnlock()
	return len(fake.hijackArgsForCall)
}

func (fake *FakeHijackStreamer) HijackCalls(stub func(context.Context, string, io.Reader, rata.Params, url.Values, string) (net.Conn, *bufio.Reader, error)) {
	fake.hijackMutex.Lock()
	defer fake.hijackMutex.Unlock()
	fake.HijackStub = stub
}

func (fake *FakeHijackStreamer) HijackArgsForCall(i int) (context.Context, string, io.Reader, rata.Params, url.Values, string) {
	fake.hijackMutex.RLock()
	defer fake.hijackMutex.RUnlock()
	argsForCall := fake.hijackArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6
}

func (fake *FakeHijackStreamer) HijackReturns(result1 net.Conn, result2 *bufio.Reader, result3 error) {
	fake.hijackMutex.Lock()
	defer fake.hijackMutex.Unlock()
	fake.HijackStub = nil
	fake.hijackReturns = struct {
		result1 net.Conn
		result2 *bufio.Reader
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeHijackStreamer) HijackReturnsOnCall(i int, result1 net.Conn, result2 *bufio.Reader, result3 error) {
	fake.hijackMutex.Lock()
	defer fake.hijackMutex.Unlock()
	fake.HijackStub = nil
	if fake.hijackReturnsOnCall == nil {
		fake.hijackReturnsOnCall = make(map[int]struct {
			result1 net.Conn
			result2 *bufio.Reader
			result3 error
		})
	}
	fake.hijackReturnsOnCall[i] = struct {
		result1 net.Conn
		result2 *bufio.Reader
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeHijackStreamer) Stream(arg1 string, arg2 io.Reader, arg3 rata.Params, arg4 url.Values, arg5 string) (io.ReadCloser, error) {
	fake.streamMutex.Lock()
	ret, specificReturn := fake.streamReturnsOnCall[len(fake.streamArgsForCall)]
	fake.streamArgsForCall = append(fake.streamArgsForCall, struct {
		arg1 string
		arg2 io.Reader
		arg3 rata.Params
		arg4 url.Values
		arg5 string
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.StreamStub
	fakeReturns := fake.streamReturns
	fake.recordInvocation("Stream", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.streamMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeHijackStreamer) StreamCallCount() int {
	fake.streamMutex.RLock()
	defer fake.streamMutex.RUnlock()
	return len(fake.streamArgsForCall)
}

func (fake *FakeHijackStreamer) StreamCalls(stub func(string, io.Reader, rata.Params, url.Values, string) (io.ReadCloser, error)) {
	fake.streamMutex.Lock()
	defer fake.streamMutex.Unlock()
	fake.StreamStub = stub
}

func (fake *FakeHijackStreamer) StreamArgsForCall(i int) (string, io.Reader, rata.Params, url.Values, string) {
	fake.streamMutex.RLock()
	defer fake.streamMutex.RUnlock()
	argsForCall := fake.streamArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeHijackStreamer) StreamReturns(result1 io.ReadCloser, result2 error) {
	fake.streamMutex.Lock()
	defer fake.streamMutex.Unlock()
	fake.StreamStub = nil
	fake.streamReturns = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeHijackStreamer) StreamReturnsOnCall(i int, result1 io.ReadCloser, result2 error) {
	fake.streamMutex.Lock()
	defer fake.streamMutex.Unlock()
	fake.StreamStub = nil
	if fake.streamReturnsOnCall == nil {
		fake.streamReturnsOnCall = make(map[int]struct {
			result1 io.ReadCloser
			result2 error
		})
	}
	fake.streamReturnsOnCall[i] = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeHijackStreamer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.hijackMutex.RLock()
	defer fake.hijackMutex.RUnlock()
	fake.streamMutex.RLock()
	defer fake.streamMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeHijackStreamer) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ connection.HijackStreamer = new(FakeHijackStreamer)
