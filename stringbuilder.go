package golibs

import (
	"bytes"
)

type StringBuilder struct {
	buf bytes.Buffer
}

func NewStringBuilder() *StringBuilder {
	return &StringBuilder{buf: bytes.Buffer{}}
}

func (this *StringBuilder) Append(obj string) *StringBuilder {
	this.buf.WriteString(obj)
	return this
}

func (this *StringBuilder) ToString() string {
	return this.buf.String()
}
