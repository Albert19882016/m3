// Code generated by "esc -modtime 12345 -prefix openapi/ -pkg openapi -ignore .go -o openapi/assets.go ."; DO NOT EDIT.

// Copyright (c) 2018 Uber Technologies, Inc.
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

package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/asset-gen.sh": {
		local:   "asset-gen.sh",
		size:    238,
		modtime: 12345,
		compressed: `
H4sIAAAAAAAC/0zKz0rEMBDH8Xue4rfTnBbSsP45LR5EfAHrTUTWdpIO0hlJIgjiu0sr6M5l4PP9dbv4
Khrr7NztMNw/vgwPdzf+4JIVCERB/s8p7o+YzAGAJOzwhDDBC56PaDPrFtYbTZvoB2+QxG2fx9lAmZXL
qYlmpGILvNBvrSPCYlOThUGHi8ura0J4L5zkE+S/pOv28Tuu9pb/gRAkqxVGnw3BzqanWrnVPhuBenKT
KbufAAAA//9BiTev7gAAAA==
`,
	},

	"/index.html": {
		local:   "openapi/index.html",
		size:    636,
		modtime: 12345,
		compressed: `
H4sIAAAAAAAC/0ySQW/bMAyF7/kVjC+9RJaHDtiQyd6wpceuQ9DLblUk2lYrS55IpzC2/ffBUdLlRr4n
fXwgqNa7h2+PP3/cQc+Db1ZqLcTq+8Pj3Rb2U4CnQb8gaCJk0WEQvyZM8xO4FuY4QTbDDKbXoUMCjsC9
I2idx/VKiGalMhZA9ajtUgAoduyxub/dfYU97qJRMivZHZD1QkyEXBcTt+JjIa+9oAesi6PD1zEmLsDE
wBi4Ll6d5b62eHQGxanZgAuOnfaCjPZYvyvOIO/CC/QJ27romUfaStnGwFR2MXYe9eioNHGQhuhzqwfn
5/p+8TElzdvbqtq8r6rNh6r6s4+HyPFaKiChrwvi2SP1iHwZelJyDXCIdobf5wZg0KlzYQvVpzdp1Na6
0F1pfzNHvoGUvKxVLbzznIQ2GqARjZiSr2/iiEGPThJrdkYuRjkP/qZR8vT0Es8kNzJQMv+XYmwon8mi
d8dUBmQZxiF/+uI1I7E8TMF6pCyWxDpY7WPA8pmKZsl6ouawOaOS+Sj+BQAA//8by2IcfAIAAA==
`,
	},

	"/spec.yml": {
		local:   "openapi/spec.yml",
		size:    18580,
		modtime: 12345,
		compressed: `
H4sIAAAAAAAC/+xcX2/bOBJ/96fgqvdw+5AoTXp7gN+cOJsaSF0jCRa4XRywtDiSuSuRKjlsmhb33Q/U
P0uWYkm2m6SB/JJYHA6HMz/Ob0hJfkPmH+8ux+TGCPJnRP8GQrUGPApAHH0yoB7+JNwnD9KQtFE8EG9F
RQCaoCS44pr4PISfRvqeBgGoMXFOj0+cERe+HI8IQY4hjInz4Wx67owIYaA9xWPkUoyJMyGMa1R8aRAY
QR4B0aA4aMIo0iXVQIzmIiAfzu5ufyd+KCn+8o54MooVaM2lOCb/kYZ4VBCfC0akQRJJBYQu7b92VEKR
/LFCjMeuG52x5XHAcWWWx1y60Zn7338+2vQzkYpIQf644vjeLFNJPXbdTMqTUdLLjc5+PrZz+wxKp/N6
e3xinUCIJwVSD60nCBE0Sl1xPiVXUgYhkCslTewkrUaFY+IUY9gGfRwkYslQvlQmct/8lP61A9t+IfdA
aKgMMImptwJynTaR09SU2gi1WbjLUC7diGoE5V7PLi7nt5fOaCU12m5SY6L/36cnb52Rjc2C4mpMHJfG
3P381hkhDfR4dLSe5/SczGkEOqYe1IN/IYXPA6PS+E7Pk36JrHY2tCxC6kEEAjtoiXNZbYVLeiZB0F2R
FS5rKvTkwKwrmGYtR/ecAfGN8GyDdkbaW0EEiWsS7zujmOJK25i5GtRn7oFOY1B4II1nABlyCEl9S7LP
UZN3K035BW2iiKqHMXGuAOtuToVkDIpaY2dsTJyi/Qowl/Ck0CaZQzEKjeOQe0k39y8tRS4aK8mM10lU
gY6l0FCa2enJyfrLpoedUkviVFqWJeQfCvwxcd64DHwueOJ+d16azk024FrRu5N3Bx7vCgQo7l0qJdVa
wb8OPq/6OLFdqQ14eRQtj2JlwhihYgMuLWiZMPZ90RJTRSNAUCXhbFEuJXtYO5GL2qW6V7djZcLYDXwy
oPFFYfXklWD1sbTnfiv+nU3/lypmEALCYXA9TXT1hnba7dnQXfLJBsgtj6wvKfhkuAI2JqgMFJfxIbZa
bJ0lgieFc+q3hFhVlEz56cH8ehL8xqIpqpOttcJRUwXVWifgCjaKqeYVUjS/jlJhUZpOPf2+BArvHsaM
wrnQSIUH6W4NOgf0NbD5ojSZ52Dz7XB6PWzehaO7Azfj6N4pKO03CcNXkIi2EefTEo0npWJcUJSqD+Nc
rLs9EvqSxDYOKisayOgFkdGeER7oaaCnl0JPe0K5Qli75KuBub4Hc9EgUBD0Ja5J0esRKKwFttFW5fi4
OfyRFRo46yk5a6/gFmeiNrY9easa64G8BvI6GHnthekKdXXNWYuBt57qaM+1PcYHPhmaWStoyL/usMm2
fV9P7rKzGZLXMx4kdMX3nhV6DfG7VOkD9AfoH24n0hX5e9F7Dfe9KX4A/QD6A5Qx3/L9Uo+7+r3vGOR3
9Yu9ma9k1OtU8Znv86+d9GPd5h/KmZ1xfpijxwT5ZVXDEhiWwPOUNb1XwCHOL+rHcl1xn8xiMYB/AH//
8iZ/Etz1FNAc4W3HNNPS4+Nbi5r8cXTQaVVzz3FFIqmRyGRSmjDwqQkRWDO0c/MuEut++CJ+WpnOc1Tx
mxa8VqCXWmzfhmeiU5VZwpHLv8DLIhEri0Hk60AkKNieozI8r6W2Pxr9Mc5eqyiZ9rGsopNdSylRo6Lx
paDLEFjNxqWUIdAC4X5o9Kqj7L3iCPpOXsgo4ngtg7YOnv1iupqiIKZcdRZGENY5Hzt5+WZDvMhLgsZ6
JbHjqFww+NJtxFlJ1Ha/aTS4U0yLuS5AccnmVEhds5QLhABKy8mXKqKYtvzyLr++DKX39y3/CvtpMb4P
6leDRh1C0YJq3H9WNo9dfom5emgL44b4xEdQc4kTzwOt93TyrAaRTjGGbgDcP3xNr2L0wmLANZZdvD2r
3WTylaFvKko659v05a7apMsd7YcyllhBw0VNTa80XL931sPgdHuwNaBN9WWPETae/mgvY3MX5e9z9kLP
2WnF5B525uX/d4rcLFNfohFby/1KPZSqbY7CRLcrqljrUuI6kWtfoZ5B+RnUHW+oDzpmM64/cFuytA8W
0S+JWbeAM7ZtuNxJfcLGWuobrmWYrIrkZd8W4a9StNVL98CDFbY5DQSLJRfYokw3R5UqRctlO0LUjrDE
xRXFrf62n+LN4u2WxlJhl9D1r1F/7AgeyH1J8PZz2ob9GilCS85NVyWWNqhaGuXBrC0kWfrYqxqyOnx/
ZxVr2ytuK9sJwkRp4xFxZvPZ3WxyPft9Nr9y8ouT3yaz68n59WVx5fpy8lsm0fCQ1UH4ZKfVvbG+KvaV
bqy9KAN7UVcTj5eI0g7SjSybFJW32H3K3bV8I8QaT0V2KRfn7RkkubhdpIz4DM6h9GhY+XkCLzQadygt
vmt8KvvHxqKkOtNig9GS4M5zuXK6PhDg38sU5edVWzoFHtunCF9i8BDYbfKLLBZpCR3Zjd97adQuKfO9
PCwxU8YUaL0n9T0jxR/exc2nlLukhK5b14ZT/95broqO/wcAAP//MmsCPJRIAAA=
`,
	},

	"/": {
		isDir: true,
		local: "openapi",
	},
}
