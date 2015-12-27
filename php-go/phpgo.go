package php

import (
	"fmt"
	"reflect"
)

type PHPExports struct {
	exports map[string]*PHPExport
	c       *CPHPExports
}

type PHPExport struct {
	name  string
	value reflect.Value
	ins   []*PHPType
	outs  []*PHPType
}

type PHPType struct {
	kind string
}

var modules = map[string]*PHPExports{}

var supportedArgKinds = map[reflect.Kind]struct{}{
	reflect.Int: struct{}{},
}

func Export(name string, exports map[string]interface{}) *PHPExports {
	pe, err := newPHPExports(exports)
	if err != nil {
		panic(err)
	}
	modules[name] = pe
	return pe
}

func newPHPExports(exports map[string]interface{}) (*PHPExports, error) {

	phpExports := &PHPExports{
		exports: make(map[string]*PHPExport, len(exports)),
	}

	for name, e := range exports {
		if pe, err := newPHPExport(name, e); err != nil {
			return nil, fmt.Errorf("Failed exporting `%s`: %s", name, err)
		} else {
			phpExports.exports[name] = pe
		}
	}

	phpExports.c = marshalExports(phpExports)

	return phpExports, nil
}

func newPHPType(t reflect.Type) (*PHPType, error) {

	if _, ok := supportedArgKinds[t.Kind()]; !ok {
		return nil, fmt.Errorf("Arguments and return values of kind %s are not yet supported", t.Kind())
	}

	pt := &PHPType{
		kind: t.Kind().String(),
	}

	return pt, nil
}

func (e *PHPExport) addIn(in reflect.Type) error {
	if pt, err := newPHPType(in); err != nil {
		return err
	} else {
		e.ins = append(e.ins, pt)
		return nil
	}
}

func (e *PHPExport) addOut(out reflect.Type) error {
	if pt, err := newPHPType(out); err != nil {
		return err
	} else {
		e.outs = append(e.outs, pt)
		return nil
	}
}

func newPHPExport(name string, export interface{}) (*PHPExport, error) {

	t := reflect.TypeOf(export)
	pe := &PHPExport{
		name:  name,
		value: reflect.ValueOf(export),
	}

	switch t.Kind() {
	case reflect.Func:
		for i := 0; i < t.NumIn(); i++ {
			if err := pe.addIn(t.In(i)); err != nil {
				return nil, err
			}
		}
		for i := 0; i < t.NumOut(); i++ {
			if err := pe.addOut(t.In(i)); err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("Exporting values of kind `%s` is not yet supported", t.Kind())
	}

	return pe, nil
}