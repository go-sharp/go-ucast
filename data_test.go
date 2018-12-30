package ucast

import (
	"reflect"
	"testing"

	_ "github.com/vivint/infectious"
)

func Test_ucastHelloMessage_toNetByteOrder(t *testing.T) {
	type fields struct {
		ucastMessage  ucastMessage
		isFecMsg      bool
		fecRequired   uint8
		fecPieces     uint8
		fecInterleave uint8
		stripeSize    uint16
		name          string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Seriallize a fec hello message",
			fields: fields{
				ucastMessage:  ucastMessage{msgID: 0xfffeacca},
				fecRequired:   0x14,
				fecPieces:     0x10,
				fecInterleave: 0xab,
				isFecMsg:      true,
				stripeSize:    0xfffe,
				name:          "hello",
			},
			want: []byte{
				0x1, 0xff, 0xfe, 0xac, 0xca,
				0x80, 0x14, 0x10, 0xab, 0xff, 0xfe,
				'h', 'e', 'l', 'l', 'o',
			},
			wantErr: false,
		},
		{
			name: "Seriallize a non fec hello message",
			fields: fields{
				ucastMessage: ucastMessage{msgID: 0xfffeacca},
				isFecMsg:     false,
				stripeSize:   0xfffe,
				name:         "hello",
			},
			want: []byte{
				0x1, 0xff, 0xfe, 0xac, 0xca,
				0, 0, 0, 0, 0xff, 0xfe,
				'h', 'e', 'l', 'l', 'o',
			},
			wantErr: false,
		},
		{
			name: "Seriallize a hello message with a name too long fails",
			fields: fields{
				ucastMessage: ucastMessage{msgID: 0xfffeacca},
				isFecMsg:     false,
				stripeSize:   0xfffe,
				name: "This message is definitive too long to send in a ucast hello message." +
					"This message is definitive too long to send in a ucast hello message." +
					"This message is definitive too long to send in a ucast hello message." +
					"This message is definitive too long to send in a ucast hello message." +
					"This message is definitive too long to send in a ucast hello message." +
					"This message is definitive too long to send in a ucast hello message." +
					"This message is definitive too long to send in a ucast hello message." +
					"This message is definitive too long to send in a ucast hello message.",
			},
			want:    []byte{},
			wantErr: true,
		},
		{
			name: "Seriallize a hello message with an empty name",
			fields: fields{
				ucastMessage: ucastMessage{msgID: 0xfffeacca},
				isFecMsg:     false,
				stripeSize:   0xfffe,
				name:         "",
			},
			want: []byte{
				0x1, 0xff, 0xfe, 0xac, 0xca,
				0, 0, 0, 0, 0xff, 0xfe,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := ucastHelloMessage{
				ucastMessage:  tt.fields.ucastMessage,
				isFecMsg:      tt.fields.isFecMsg,
				fecRequired:   tt.fields.fecRequired,
				fecPieces:     tt.fields.fecPieces,
				fecInterleave: tt.fields.fecInterleave,
				stripeSize:    tt.fields.stripeSize,
				name:          tt.fields.name,
			}
			got, err := u.toNetByteOrder()
			if (err != nil) != tt.wantErr {
				t.Errorf("ucastHelloMessage.toNetByteOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("ucastHelloMessage.toNetByteOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ucastHelloMessage_fromNetByteOrder(t *testing.T) {
	type fields struct {
		ucastMessage  ucastMessage
		isFecMsg      bool
		fecRequired   uint8
		fecPieces     uint8
		fecInterleave uint8
		stripeSize    uint16
		name          string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &ucastHelloMessage{
				ucastMessage:  tt.fields.ucastMessage,
				isFecMsg:      tt.fields.isFecMsg,
				fecRequired:   tt.fields.fecRequired,
				fecPieces:     tt.fields.fecPieces,
				fecInterleave: tt.fields.fecInterleave,
				stripeSize:    tt.fields.stripeSize,
				name:          tt.fields.name,
			}
			if err := u.fromNetByteOrder(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ucastHelloMessage.fromNetByteOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
