package auto

import (
	"github.com/craiggwilson/go-mapper/pkg/auto/naming"
	"github.com/craiggwilson/go-mapper/pkg/core"
	"github.com/craiggwilson/go-mapper/pkg/reflecth"
)

type withAccessorOpt interface {
	WithAccessor(reflecth.Accessor)
}

type withConverterOpt interface {
	WithConverter(reflecth.Converter)
}

type withConverterFactoryOpt interface {
	WithConverterFactory(reflecth.ConverterFactory)
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

func WithFieldAccessor(a reflecth.Accessor) func(fieldOpts) {
	return func(opt fieldOpts) {
		opt.WithAccessor(a)
	}
}

func WithFieldConverter(c reflecth.Converter) func(fieldOpts) {
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

func WithStructConverterFactory(cf reflecth.ConverterFactory) func(structOpts) {
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
