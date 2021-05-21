package params

var (
	validateErrTypeParse = "%v can not parse to %s"
	validateUrlParse     = "%s is not a valid url: %s"
)

type paramValidateFunc func(p interface{}) error
type ParamValidateItem struct {
	Value        string
	ValidateFunc paramValidateFunc
}

func NewParamValidateItem(value string, validateFunc paramValidateFunc) *ParamValidateItem {
	return &ParamValidateItem{Value: value, ValidateFunc: validateFunc}
}

type ParamValidator struct {
	items []ParamValidateItem
}

func NewParamValidator(rules []ParamValidateItem) *ParamValidator {
	return &ParamValidator{items: rules}
}

func (p *ParamValidator) AppendRule(validateItem ParamValidateItem) {
	p.items = append(p.items, validateItem)
}

func (p *ParamValidator) Check() error {
	for i := range p.items {
		if err := p.items[i].ValidateFunc(p.items[i].Value); err != nil {
			return err
		}
	}
	return nil
}
