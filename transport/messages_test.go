package transport

import (
	"bytes"
	"reflect"
	"testing"
)

func TestHeaderSerialisation(t *testing.T) {
	testCases := []struct {
		name   string
		header messageHeader
		want   messageHeader
	}{
		{
			name: "Header message",
			header: messageHeader{
				message: message{
					msgID:   0x112233aabbcc,
					msgType: 0xae,
					flags:   0,
				},
				stripeL: 0xff11,
				fecTot:  0xa9,
				fecRec:  0xb2,
				fecInt:  0xc4,
			},
			want: messageHeader{
				message: message{
					msgID:   0x112233aabbcc,
					msgType: 0xae,
					flags:   0,
				},
				stripeL: 0xff11,
			},
		},
		{
			name: "Fec header message",
			header: messageHeader{
				message: message{
					msgID:   0x112233aabbcc,
					msgType: 0xae,
					flags:   fecMsgFlag | 0xfe11,
				},
				stripeL: 0xff11,
				fecTot:  0xa9,
				fecRec:  0xb2,
				fecInt:  0xc4,
			},
			want: messageHeader{
				message: message{
					msgID:   0x112233aabbcc,
					msgType: 0xae,
					flags:   fecMsgFlag | 0xfe11,
				},
				stripeL: 0xff11,
				fecTot:  0xa9,
				fecRec:  0xb2,
				fecInt:  0xc4,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			var result messageHeader
			// act
			tc.header.toBytes(&buf)
			result.fromBytes(buf.Bytes())

			// assert
			if !reflect.DeepEqual(tc.want, result) {
				t.Errorf("messageHeader serialisation: got = %v, want %v", result, tc.want)
			}
		})
	}
}

func TestDataSerialisation(t *testing.T) {
	testCases := []struct {
		name   string
		header messageData
		want   messageData
	}{
		{
			name: "Data message",
			header: messageData{
				message: message{
					msgID:   0x112233aabbcc,
					msgType: 0xae,
					flags:   0,
				},
				data:     []byte{123, 222, 23, 56, 11},
				stripeNr: 0xa9,
				fecPad:   0xb2,
				fecNr:    0xc4,
			},
			want: messageData{
				message: message{
					msgID:   0x112233aabbcc,
					msgType: 0xae,
					flags:   0,
				},
				data:     []byte{123, 222, 23, 56, 11},
				stripeNr: 0xa9,
			},
		},
		{
			name: "Fec data message",
			header: messageData{
				message: message{
					msgID:   0x112233aabbcc,
					msgType: 0xae,
					flags:   fecMsgFlag | 0xfe11,
				},
				data:     []byte{123, 222, 23, 56, 11},
				stripeNr: 0xa9,
				fecPad:   0xb2,
				fecNr:    0xc4,
			},
			want: messageData{
				message: message{
					msgID:   0x112233aabbcc,
					msgType: 0xae,
					flags:   fecMsgFlag | 0xfe11,
				},
				data:     []byte{123, 222, 23, 56, 11},
				stripeNr: 0xa9,
				fecPad:   0xb2,
				fecNr:    0xc4,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			var result messageData
			// act
			tc.header.toBytes(&buf)
			result.fromBytes(buf.Bytes())

			// assert
			if !reflect.DeepEqual(tc.want, result) {
				t.Errorf("messageData serialisation: got = %v, want %v", result, tc.want)
			}
		})
	}
}

func Test_flags_isSet(t *testing.T) {
	type args struct {
		fl []flags
	}
	tests := []struct {
		name string
		f    flags
		args args
		want bool
	}{
		{
			name: "Single flag is set",
			f:    1 << 6,
			args: args{fl: []flags{1 << 6}},
			want: true,
		},
		{
			name: "Multiple flag is set",
			f:    1<<6 | 1<<9 | 1<<2,
			args: args{fl: []flags{1 << 2, 1 << 6}},
			want: true,
		},
		{
			name: "Multiple flag is not set",
			f:    1<<6 | 1<<9 | 1<<2,
			args: args{fl: []flags{1 << 9, 1 << 5}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.isSet(tt.args.fl...); got != tt.want {
				t.Errorf("flags.isSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
