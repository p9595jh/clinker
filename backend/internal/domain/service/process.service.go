package service

import (
	"clinker-backend/common/util"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/p9595jh/transform"
)

type ProcessService struct {
	validator   *validator.Validate
	transformer transform.Transformer
}

func NewProcessService(validator *validator.Validate, transformer transform.Transformer) *ProcessService {
	return &ProcessService{
		validator:   validator,
		transformer: transformer,
	}
}

func (s *ProcessService) Initializer() {
	regexps := map[string]*regexp.Regexp{
		"date":    regexp.MustCompile(`\d{4}-\d{2}-\d{2}`),
		"ethAddr": regexp.MustCompile("^(0x)?[a-fA-F0-9]{40,40}$"),
		"txHash":  regexp.MustCompile("^(0x)?[a-fA-F0-9]{64,64}$"),
	}
	for tag, regex := range regexps {
		s.validator.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
			return regex.MatchString(fl.Field().String())
		})
	}

	s.transformer.RegisterTransformer("hex", transform.F1(func(s1, s2 string) string {
		s1 = strings.TrimPrefix(s1, "0x")
		s1 = strings.ToLower(s1)
		return s1
	}))
	s.transformer.RegisterTransformer("timeFormat", transform.F2(func(t time.Time, s string) string {
		return t.Format(util.DateFormat)
	}))
	s.transformer.RegisterTransformer("add0x", transform.F1(func(s1, s2 string) string {
		return "0x" + s1
	}))
}

// target must be pointer
func (s *ProcessService) Decode(params map[string]string, target any) error {
	return mapstructure.Decode(params, &target)
}

func (s *ProcessService) Validate(data any) error {
	return s.validator.Struct(data)
}

func (s *ProcessService) Transform(src, dst any) error {
	if dst == nil {
		return s.transformer.Transform(src)
	} else {
		return s.transformer.Mapping(src, dst)
	}
}

func (s *ProcessService) Pre(params map[string]string, v, mapped any) error {
	if err := s.Decode(params, v); err != nil {
		return err
	}

	if err := s.Validate(v); err != nil {
		return err
	}

	if mapped == nil {
		return s.transformer.Transform(v)
	} else {
		return s.transformer.Mapping(v, mapped)
	}
}

func (s *ProcessService) PreWithParser(parser func(any) error, v, mapped any) error {
	if err := parser(v); err != nil {
		return err
	}

	if err := s.Validate(v); err != nil {
		return err
	}

	if mapped == nil {
		return s.transformer.Transform(v)
	} else {
		return s.transformer.Mapping(v, mapped)
	}
}

func (s *ProcessService) Post(v, mapped any) error {
	if mapped == nil {
		return s.transformer.Transform(v)
	} else {
		return s.transformer.Mapping(v, mapped)
	}
}
