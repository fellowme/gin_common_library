package gin_remote_service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	httpUrl "net/url"
	"time"
)

var defaultClient = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives:  true,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:        100,
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 100,
	},
}

type ResponseData struct {
	Data []byte
	Code int
}

func doHttp(ctx context.Context, method, url string, payload []byte, opt *option) (*ResponseData, error) {
	// 返回 加的测试mock 数据
	if mock := opt.mock; mock != nil {
		return &ResponseData{
			Data: mock(),
			Code: http.StatusOK,
		}, nil
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(payload))
	if err != nil {
		opt.logger.Error("http.NewRequestWithContext",
			zap.String("method", method),
			zap.String("url", url),
			zap.Any("error", err),
			zap.ByteString("payload", payload))
		return &ResponseData{
			Data: nil,
			Code: -1,
		}, errors.Wrapf(err, "http.NewRequestWithContext new request [%s %s] err", method, url)
	}

	for key, value := range opt.header {
		req.Header.Set(key, value)
	}

	resp, err := defaultClient.Do(req)
	if err != nil {
		header, _ := json.Marshal(opt.header)
		opt.logger.Error("defaultClient.Do",
			zap.String("method", method),
			zap.String("url", url),
			zap.Any("error", err),
			zap.ByteString("payload", payload),
			zap.String("header", string(header)))
		err = errors.Wrapf(err, "defaultClient.Do request [%s %s] err", method, url)
		return &ResponseData{
			Data: nil,
			Code: _StatusDoReqErr,
		}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		header, _ := json.Marshal(opt.header)
		opt.logger.Error("ioutil.ReadAll",
			zap.String("method", method),
			zap.String("url", url),
			zap.Any("error", err),
			zap.ByteString("payload", payload),
			zap.String("header", string(header)))
		err = errors.Wrapf(err, "ioutil.ReadAll read resp body from [%s %s] err", method, url)
		return &ResponseData{
			Data: nil,
			Code: _StatusReadRespErr,
		}, err
	}
	defer func() {
		if closeError := resp.Body.Close(); closeError != nil {
			opt.logger.Error("resp.Body.Close",
				zap.Any("error", closeError))
		}
	}()
	return &ResponseData{
		Data: body,
		Code: resp.StatusCode,
	}, nil
}

// Get get 请求
func Get(url string, form httpUrl.Values, options ...Option) (resp *ResponseData, err error) {
	return withoutBody(http.MethodGet, url, form, options...)
}

// Delete delete 请求
func Delete(url string, form httpUrl.Values, options ...Option) (resp *ResponseData, err error) {
	return withoutBody(http.MethodDelete, url, form, options...)
}

// PostForm post form 请求
func PostForm(url string, form httpUrl.Values, options ...Option) (resp *ResponseData, err error) {
	return withFormBody(http.MethodPost, url, form, options...)
}

// PostJSON post json 请求
func PostJSON(url string, raw json.RawMessage, options ...Option) (resp *ResponseData, err error) {
	return withJSONBody(http.MethodPost, url, raw, options...)
}

// PutForm put form 请求
func PutForm(url string, form httpUrl.Values, options ...Option) (resp *ResponseData, err error) {
	return withFormBody(http.MethodPut, url, form, options...)
}

// PutJSON put json 请求
func PutJSON(url string, raw json.RawMessage, options ...Option) (resp *ResponseData, err error) {
	return withJSONBody(http.MethodPut, url, raw, options...)
}

// PatchFrom patch form 请求
func PatchFrom(url string, form httpUrl.Values, options ...Option) (resp *ResponseData, err error) {
	return withFormBody(http.MethodPatch, url, form, options...)
}

// PatchJSON patch json 请求
func PatchJSON(url string, raw json.RawMessage, options ...Option) (resp *ResponseData, err error) {
	return withJSONBody(http.MethodPatch, url, raw, options...)
}

func withJSONBody(method, url string, raw json.RawMessage, options ...Option) (resp *ResponseData, err error) {
	opt := getOption()
	for _, f := range options {
		f(opt)
	}
	if url == "" {
		opt.logger.Error("withoutBody url required",
			zap.String("url", ""),
			zap.String("method", method),
		)
		return resp, errors.New("url required")
	}

	if len(raw) <= 0 {
		return resp, errors.New("raw required")
	}
	opt.header["Content-Type"] = "application/json; charset=utf-8"
	if opt.jaegerContext != nil {
		c := opt.jaegerContext
		tracer := c.Value("Tracer").(opentracing.Tracer)
		parentSpanContext := c.Value("ParentSpanContext").(opentracing.SpanContext)
		span := opentracing.StartSpan(
			url,
			opentracing.ChildOf(parentSpanContext),
			opentracing.Tags{
				"method":              method,
				"url":                 url,
				"data":                raw,
				"header":              c.Request.Header,
				string(ext.Component): c.Request.Proto,
				"serverName":          gin_config.ServerConfigSettings.Server.ServerName,
			},
			ext.SpanKindRPCClient,
		)
		defer span.Finish()
		injectErr := tracer.(opentracing.Tracer).Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if injectErr != nil {
			zap.L().Error("injectErr error", zap.Any("error", injectErr))
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), opt.ttl)
	defer cancel()
	retryTimes := opt.retryTimes
	for k := 0; k < retryTimes; k++ {
		resp, err = doHttp(ctx, method, url, raw, opt)
		if resp == nil || shouldRetry(ctx, resp.Code) {
			if err == nil {
				err = errors.New("resp doHttp返回错误")
			}
			if resp == nil {
				resp = &ResponseData{
					Data: nil,
					Code: 0,
				}
			}
			message := fmt.Sprintf("withoutBody doHttp 第 %d 次失败", k)
			opt.logger.Error(message,
				zap.String("url", url),
				zap.String("method", method),
				zap.String("err", err.Error()),
				zap.Int("code", resp.Code),
				zap.ByteString("data", resp.Data),
			)
			time.Sleep(opt.retryDelay)
			continue
		}
		return
	}
	return
}

func withFormBody(method, url string, form httpUrl.Values, options ...Option) (resp *ResponseData, err error) {
	opt := getOption()
	for _, f := range options {
		f(opt)
	}
	if url == "" {
		opt.logger.Error("withoutBody url required",
			zap.String("url", ""),
			zap.String("method", method),
		)
		return resp, errors.New("url required")
	}

	if len(form) <= 0 {
		return resp, errors.New("form required")
	}
	opt.header["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	ctx, cancel := context.WithTimeout(context.Background(), opt.ttl)
	defer cancel()
	retryTimes := opt.retryTimes
	formValue := form.Encode()
	if opt.jaegerContext != nil {
		c := opt.jaegerContext
		tracer := c.Value("Tracer").(opentracing.Tracer)
		parentSpanContext := c.Value("ParentSpanContext").(opentracing.SpanContext)
		span := opentracing.StartSpan(
			url,
			opentracing.ChildOf(parentSpanContext),
			opentracing.Tags{
				"method":              method,
				"url":                 url,
				"data":                formValue,
				"header":              c.Request.Header,
				string(ext.Component): c.Request.Proto,
				"serverName":          gin_config.ServerConfigSettings.Server.ServerName,
			},
			ext.SpanKindRPCClient,
		)
		defer span.Finish()
		injectErr := tracer.(opentracing.Tracer).Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if injectErr != nil {
			zap.L().Error("injectErr error", zap.Any("error", injectErr))
		}
	}
	for k := 0; k < retryTimes; k++ {
		resp, err = doHttp(ctx, method, url, []byte(formValue), opt)
		if resp == nil || shouldRetry(ctx, resp.Code) {
			if err == nil {
				err = errors.New("resp doHttp返回错误")
			}
			if resp == nil {
				resp = &ResponseData{
					Data: nil,
					Code: 0,
				}
			}
			message := fmt.Sprintf("withoutBody doHttp 第 %d 次失败", k)
			opt.logger.Error(message,
				zap.String("url", url),
				zap.String("method", method),
				zap.String("err", err.Error()),
				zap.Int("code", resp.Code),
				zap.ByteString("data", resp.Data),
			)
			time.Sleep(opt.retryDelay)
			continue
		}
		return
	}
	return
}

func withoutBody(method, url string, form httpUrl.Values, options ...Option) (resp *ResponseData, err error) {
	opt := getOption()
	for _, f := range options {
		f(opt)
	}
	if url == "" {
		opt.logger.Error("withoutBody url required",
			zap.String("url", ""),
			zap.String("method", method),
		)
		return resp, errors.New("url required")
	}

	if len(form) > 0 {
		if url, err = addFormValuesIntoURL(url, form); err != nil {
			opt.logger.Error("withoutBody addFormValuesIntoURL",
				zap.String("url", url),
				zap.String("method", method),
				zap.Any("err", err),
				zap.Any("form", form),
			)
			return
		}
	}
	opt.header["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	if opt.jaegerContext != nil {
		c := opt.jaegerContext
		tracer := c.Value("Tracer").(opentracing.Tracer)
		parentSpanContext := c.Value("ParentSpanContext").(opentracing.SpanContext)
		span := opentracing.StartSpan(
			url,
			opentracing.ChildOf(parentSpanContext),
			opentracing.Tags{
				"method":              method,
				"url":                 url,
				"header":              c.Request.Header,
				string(ext.Component): "HTTP",
				"serverName":          gin_config.ServerConfigSettings.Server.ServerName,
			},
			ext.SpanKindRPCClient,
		)
		defer span.Finish()

		injectErr := tracer.(opentracing.Tracer).Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if injectErr != nil {
			zap.L().Error("injectErr error", zap.Any("error", injectErr))
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), opt.ttl)
	defer cancel()
	retryTimes := opt.retryTimes
	for k := 0; k < retryTimes; k++ {
		resp, err = doHttp(ctx, method, url, nil, opt)
		if resp == nil || shouldRetry(ctx, resp.Code) {
			if err == nil {
				err = errors.New("resp doHttp返回错误")
			}
			if resp == nil {
				resp = &ResponseData{
					Data: nil,
					Code: 0,
				}
			}
			message := fmt.Sprintf("withoutBody doHttp 第 %d 次失败", k)
			opt.logger.Error(message,
				zap.String("url", url),
				zap.String("method", method),
				zap.String("err", err.Error()),
				zap.Int("code", resp.Code),
				zap.ByteString("data", resp.Data),
			)
			time.Sleep(opt.retryDelay)
			continue
		}
		return
	}
	return
}
