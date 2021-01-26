package auto

import (
	"github.com/craiggwilson/go-mapper/pkg/auto/accessor"
	"github.com/craiggwilson/go-mapper/pkg/auto/converter"
	"github.com/craiggwilson/go-mapper/pkg/auto/naming"
	"github.com/craiggwilson/go-mapper/pkg/core"
)

type withAccessorOpt interface {
	WithAccessor(accessor.Accessor)
}

type withConverterOpt interface {
	WithConverter(converter.Converter)
}

type withConverterFactoryOpt interface {
	WithConverterFactory(converter.Factory)
}

type withFieldOpt interface {
	WithField(name string, opts ...func(fieldOpts))
}

type withMapperOpt interface {
	WithMapper(core.Mapper)
}

type withNamingStrategyOpt interface {
	WithNamingStrategy(naming.Strategy)
}

type fieldOpts interface {
	withAccessorOpt
	withConverterOpt
	withMapperOpt
	withNamingStrategyOpt
}

type structOpts interface {
	withConverterFactoryOpt
	withFieldOpt
	withNamingStrategyOpt
}

func WithFieldAccessor(a accessor.Accessor) func(fieldOpts) {
	return func(opt fieldOpts) {
		opt.WithAccessor(a)
	}
}

func WithFieldConverter(c converter.Converter) func(fieldOpts) {
	return func(opt fieldOpts) {
		opt.WithConverter(c)
	}
}

func WithFieldMapper(m core.Mapper) func(opt fieldOpts) {
	return func(opt fieldOpts) {
		opt.WithMapper(m)
	}
}

func WithFieldNamingConvention(ns naming.Strategy) func(fieldOpts) {
	return func(opt fieldOpts) {
		opt.WithNamingStrategy(ns)
	}
}

func WithStructConverterFactory(cf converter.Factory) func(structOpts) {
	return func(opt structOpts) {
		opt.WithConverterFactory(cf)
	}
}

func WithStructField(name string, opts ...func(fieldOpts)) func(structOpts) {
	return func(opt structOpts) {
		opt.WithField(name, opts...)
	}
}

func WithStructNamingConvention(ns naming.Strategy) func(structOpts) {
	return func(opt structOpts) {
		opt.WithNamingStrategy(ns)
	}
}
