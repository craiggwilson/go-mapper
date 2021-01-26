package auto

import (
	"github.com/craiggwilson/go-mapper/pkg/core"
	"github.com/craiggwilson/go-mapper/pkg/reflecth"
)

type withConverterOpt interface {
	WithConverter(c reflecth.Converter)
}

type withConverterFactoryOpt interface {
	WithConverterFactory(cf reflecth.ConverterFactory)
}

type withFieldOpt interface {
	WithField(name string, opts ...func(fieldOpts))
}

type withMapperOpt interface {
	WithMapper(m core.Mapper)
}

type withNamingConventionOpt interface {
	WithNamingConvention(nc NamingConvention)
}

type fieldOpts interface {
	withConverterOpt
	withMapperOpt
	withNamingConventionOpt
}

type structOpts interface {
	withConverterFactoryOpt
	withFieldOpt
	withNamingConventionOpt
}

func WithFieldConverter(c reflecth.Converter) func(fieldOpts) {
	return func(opt fieldOpts) {
		opt.WithConverter(c)
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

func WithFieldMapper(m core.Mapper) func(opt fieldOpts) {
	return func(opt fieldOpts) {
		opt.WithMapper(m)
	}
}

func WithFieldNamingConvention(nc NamingConvention) func(fieldOpts) {
	return func(opt fieldOpts) {
		opt.WithNamingConvention(nc)
	}
}

func WithStructNamingConvention(nc NamingConvention) func(structOpts) {
	return func(opt structOpts) {
		opt.WithNamingConvention(nc)
	}
}
