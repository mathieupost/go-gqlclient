package gqlclient

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Message    string
		Path       ast.Path
		Locations  []gqlerror.Location
		Extensions map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Message",
			fields: fields{Message: "message"},
			want:   "graphql: message",
		},
		{
			name: "MessageWithPath",
			fields: fields{
				Message: "message",
				Path:    []ast.PathElement{ast.PathName("users"), ast.PathIndex(1), ast.PathName("firstName")},
			},
			want: "graphql: users[1].firstName: message",
		},
		{
			name: "MessageWithLocation",
			fields: fields{
				Message: "message",
				Locations: []gqlerror.Location{{
					Line:   1,
					Column: 2,
				}},
			},
			want: "graphql:1:2: message",
		},
		{
			name: "MessageWithPathAndLocation",
			fields: fields{
				Message: "message",
				Path:    []ast.PathElement{ast.PathName("users"), ast.PathIndex(1), ast.PathName("firstName")},
				Locations: []gqlerror.Location{{
					Line:   1,
					Column: 2,
				}},
			},
			want: "graphql:1:2: users[1].firstName: message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Message:    tt.fields.Message,
				Path:       tt.fields.Path,
				Locations:  tt.fields.Locations,
				Extensions: tt.fields.Extensions,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
