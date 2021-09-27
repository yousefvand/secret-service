package service_test

import (
	"reflect"
	"regexp"
	"sync"
	"testing"

	"github.com/yousefvand/secret-service/pkg/crypto"
	"github.com/yousefvand/secret-service/pkg/service"
)

func Test_UUID(t *testing.T) {

	const uuidLength = 32
	uuidRegex, err := regexp.Compile("[a-z0-9]+")
	if err != nil {
		t.Errorf("UUID regex compilation failed. Error: %s", err.Error())
	}

	for i := 0; i < 100; i++ {
		uuid := service.UUID()
		if len(uuid) != uuidLength {
			t.Errorf("UUID() = %v, want %v", uuid, uuidLength)
		}
		if !uuidRegex.MatchString(uuid) {
			t.Errorf("UUID has illegal character. Allowed [a-z] and [0-9]. Got: %s", uuid)
		}
	}
}

func Test_Path2Name(t *testing.T) {
	type args struct {
		path string
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "Path to name",
			args: args{
				path: "/a/b/c/xyz",
				name: "Foo",
			},
			want:  "a.b.c.Foo",
			want1: "xyz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := service.Path2Name(tt.args.path, tt.args.name)
			if got != tt.want {
				t.Errorf("Path2Name() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Path2Name() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestIsMapSubsetSingleMatch(t *testing.T) {

	type args struct {
		mapSet    map[string]string
		mapSubset map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty subset",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d"},
				mapSubset: map[string]string{},
			},
			want: true,
		},
		{
			name: "map[string]string subset bigger",
			args: args{
				mapSet:    map[string]string{"a": "b", "e": "f"},
				mapSubset: map[string]string{"a": "b", "c": "d", "e": "f"},
			},
			want: false,
		},
		{
			name: "map[string]string subset",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: map[string]string{"a": "b", "e": "f"},
			},
			want: true,
		},
		{
			name: "map[string]string subset mismatch value",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: map[string]string{"a": "b", "e": "f", "c": "z"},
			},
			want: true,
		},
		{
			name: "map[string]string subset mismatch key",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: map[string]string{"a": "b", "e": "f", "z": "d"},
			},
			want: true,
		},
		{
			name: "map[string]string subset complete mismatch",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d"},
				mapSubset: map[string]string{"e": "f"},
			},
			want: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			lock := new(sync.RWMutex)
			if got := service.IsMapSubsetSingleMatch(tt.args.mapSet,
				tt.args.mapSubset, lock); got != tt.want {
				t.Errorf("IsMapSubset() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestIsMapSubsetFullMatch(t *testing.T) {

	type args struct {
		mapSet    map[string]string
		mapSubset map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty subset",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d"},
				mapSubset: map[string]string{},
			},
			want: true,
		},
		{
			name: "map[string]string subset bigger",
			args: args{
				mapSet:    map[string]string{"a": "b", "e": "f"},
				mapSubset: map[string]string{"a": "b", "c": "d", "e": "f"},
			},
			want: false,
		},
		{
			name: "map[string]string subset",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: map[string]string{"a": "b", "e": "f"},
			},
			want: true,
		},
		{
			name: "map[string]string subset mismatch value",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: map[string]string{"a": "b", "e": "f", "c": "z"},
			},
			want: false,
		},
		{
			name: "map[string]string subset mismatch key",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: map[string]string{"a": "b", "e": "f", "z": "d"},
			},
			want: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			lock := new(sync.RWMutex)
			if got := service.IsMapSubsetFullMatch(tt.args.mapSet,
				tt.args.mapSubset, lock); got != tt.want {
				t.Errorf("IsMapSubset() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestIsMapSubsetFullMatchGeneric(t *testing.T) {

	type args struct {
		mapSet    interface{}
		mapSubset interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "complete mismatch types",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: 13,
			},
			want: false,
		},
		{
			name: "mismatch types",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: map[int]int{1: 2, 3: 4},
			},
			want: false,
		},
		{
			name: "empty subset",
			args: args{
				mapSet:    map[int]int{1: 2, 3: 4},
				mapSubset: map[int]int{},
			},
			want: true,
		},
		{
			name: "map[string]string subset bigger",
			args: args{
				mapSet:    map[string]string{"a": "b", "e": "f"},
				mapSubset: map[string]string{"a": "b", "c": "d", "e": "f"},
			},
			want: false,
		},
		{
			name: "map[string]string subset",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: map[string]string{"a": "b", "e": "f"},
			},
			want: true,
		},
		{
			name: "map[string]string subset mismatch value",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: map[string]string{"a": "b", "e": "f", "c": "z"},
			},
			want: false,
		},
		{
			name: "map[string]string subset mismatch key",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: map[string]string{"a": "b", "e": "f", "z": "d"},
			},
			want: false,
		},
		{
			name: "map[string]string subset identical",
			args: args{
				mapSet:    map[string]string{"a": "b", "c": "d", "e": "f"},
				mapSubset: map[string]string{"a": "b", "e": "f", "c": "d"},
			},
			want: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			lock := new(sync.RWMutex)
			if got := service.IsMapSubsetFullMatchGeneric(tt.args.mapSet,
				tt.args.mapSubset, lock); got != tt.want {
				t.Errorf("IsMapSubset() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestCommandExists(t *testing.T) {
	type args struct {
		cmdName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "command exists",
			args: args{cmdName: "ls"},
			want: true,
		},
		{
			name: "command exists",
			args: args{cmdName: "remisa"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.CommandExists(tt.args.cmdName); got != tt.want {
				t.Errorf("CommandExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPKCS7Padding(t *testing.T) {
	type args struct {
		plainUnpaddedData []byte
		blockSize         int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "smaller",
			args: args{
				plainUnpaddedData: []byte("aaaaaaaaaaaa"),
				blockSize:         16,
			},
			want: []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 4, 4, 4, 4},
		},
		{
			name: "same size",
			args: args{
				plainUnpaddedData: []byte("aaaaaaaaaaaaaaaa"),
				blockSize:         16,
			},
			want: []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97,
				16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16},
		},
		{
			name: "bigger",
			args: args{
				plainUnpaddedData: []byte("aaaaaaaaaaaaaaaaaaaaaaaa"),
				blockSize:         16,
			},
			want: []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97,
				97, 97, 97, 97, 97, 97, 97, 97, 8, 8, 8, 8, 8, 8, 8, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := crypto.PKCS7Padding(tt.args.plainUnpaddedData, tt.args.blockSize)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PKCS7Padding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPKCS7UnPadding(t *testing.T) {
	type args struct {
		plainPaddedData []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "smaller",
			args: args{
				plainPaddedData: []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 4, 4, 4, 4},
			},
			want: []byte("aaaaaaaaaaaa"),
		},
		{
			name: "same size",
			args: args{
				plainPaddedData: []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97,
					16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16},
			},
			want: []byte("aaaaaaaaaaaaaaaa"),
		},
		{
			name: "bigger",
			args: args{
				plainPaddedData: []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97,
					97, 97, 97, 97, 97, 97, 97, 97, 8, 8, 8, 8, 8, 8, 8, 8},
			},
			want: []byte("aaaaaaaaaaaaaaaaaaaaaaaa"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := crypto.PKCS7UnPadding(tt.args.plainPaddedData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PKCS7UnPadding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAesCBCEncrypt(t *testing.T) {

	secret := []byte("Victoria")
	key := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	iv, cipherData, err := crypto.AesCBCEncrypt(secret, key)

	if err != nil {
		t.Errorf("Encryption failed. Error: %v", err)
	}

	data, err := crypto.AesCBCDecrypt(iv, cipherData, key)

	if err != nil {
		t.Errorf("Decryption failed. Error: %v", err)
	}

	if string(data) != string(secret) {
		t.Errorf("Expected decrypted data to be: %v, got: %v", secret, data)
	}

}

func TestEncryptionAndDecryption(t *testing.T) {
	type args struct {
		key    string
		secret string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test encryption/decryption",
			args: args{
				key:    "0A1B2C3D4E5F6G7H8I9J0K1L2M3N4O5P",
				secret: "Victoria",
			},
			want: "Victoria",
		},
		{
			name: "test encryption/decryption",
			args: args{
				key:    "012345678901234567890123456789ab",
				secret: "Victoria2",
			},
			want: "Victoria2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cipher, err := crypto.EncryptAESCBC256(tt.args.key, tt.args.secret)

			if err != nil {
				t.Errorf("Encryption failed. Error: %v", err)
			}

			got, err := crypto.DecryptAESCBC256(tt.args.key, cipher)

			if err != nil {
				t.Errorf("Decryption failed. Error: %v", err)
			}

			if tt.want != got {
				t.Errorf("Expected: %s, got: %s", tt.want, got)
			}

		})
	}
}
