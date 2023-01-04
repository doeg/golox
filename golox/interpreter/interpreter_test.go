package interpreter

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/doeg/golox/golox/ast"
	"github.com/doeg/golox/golox/parser"
	"github.com/doeg/golox/golox/scanner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		input    any
		expected bool
	}{
		{
			input:    false,
			expected: false,
		},
		{
			input:    nil,
			expected: false,
		},
		{
			input:    true,
			expected: true,
		},
		{
			input:    "hello",
			expected: true,
		},
		{
			input:    0,
			expected: true,
		},
		{
			input:    1,
			expected: true,
		},
		{
			input:    -1,
			expected: true,
		},
		{
			input:    "",
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%+v", tt.input), func(t *testing.T) {
			t.Parallel()

			var output bytes.Buffer
			i := New(&output)
			result := i.isTruthy(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVisitBinaryExpression(t *testing.T) {
	tests := []struct {
		input         string
		expected      any
		expectedError error
	}{
		//
		// token.BANG_EQUAL
		//
		{
			input:    "1 != 2",
			expected: true,
		},
		{
			input:    "1 != 2 != 3",
			expected: true,
		},
		{
			input:    "1 != \"hello\"",
			expected: true,
		},
		{
			input:    "\"hello\" != false",
			expected: true,
		},
		{
			input:    "nil != nil",
			expected: false,
		},
		//
		// token.EQUAL_EQUAL
		//
		{
			input:    "nil == nil",
			expected: true,
		},
		{
			input:    "1 == 1",
			expected: true,
		},
		{
			input:    "\"hello\" == \"hello\"",
			expected: true,
		},
		{
			input:    "false == false",
			expected: true,
		},
		{
			input:    "1 == 1 == 1",
			expected: false,
		},
		{
			input:    "1 == 2",
			expected: false,
		},
		{
			input:    "1 == 2 == 3",
			expected: false,
		},
		{
			input:    "1 == \"hello\"",
			expected: false,
		},
		{
			input:    "\"hello\" == false",
			expected: false,
		},
		//
		// token.GREATER
		//
		{
			input:    "2 > 1",
			expected: true,
		},
		{
			input:    "1 > 1",
			expected: false,
		},
		{
			input:    "1 > 2",
			expected: false,
		},
		{
			input:         "1 > \"hello\"",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
		{
			input:         "\"hello\" > false",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
		//
		// token.GREATER_EQUAL
		//
		{
			input:    "1 >= 1",
			expected: true,
		}, {
			input:    "2 >= 1",
			expected: true,
		},
		{
			input:    "1 >= 2",
			expected: false,
		},
		{
			input:         "1 > \"hello\"",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
		{
			input:         "\"hello\" > false",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
		//
		// token.LESS
		//
		{
			input:    "1 < 2",
			expected: true,
		},
		{
			input:    "2 < 1",
			expected: false,
		},
		{
			input:    "1 < 1",
			expected: false,
		},
		{
			input:         "1 < \"hello\"",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
		{
			input:         "\"hello\" < false",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
		//
		// token.LESS_EQUAL
		//
		{
			input:    "1 <= 2",
			expected: true,
		},
		{
			input:    "1 <= 1",
			expected: true,
		},
		{
			input:    "2 <= 1",
			expected: false,
		},
		{
			input:         "1 <= \"hello\"",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
		{
			input:         "\"hello\" <= false",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
		//
		// token.MINUS
		//
		{
			input:    "1 - 2",
			expected: float64(-1),
		},
		{
			input:    "1 - 2 - 3",
			expected: float64(-4),
		},
		{
			input:         "1 - \"hello\"",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
		{
			input:         "\"hello\" - false",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
		//
		// token.PLUS
		//
		{
			input:    "1 + 2",
			expected: float64(3),
		},
		{
			input:    "1.2 + 3.4",
			expected: float64(4.6),
		},
		{
			input:    "0 + 0",
			expected: float64(0),
		},
		{
			input:    "\"hello\" + \"world\"",
			expected: "helloworld",
		},
		{
			input:         "\"hello\" + 1",
			expectedError: errors.New("operators must be strings or numbers"),
		},
		{
			input:         "1 + false",
			expectedError: errors.New("operators must be strings or numbers"),
		},
		{
			input:         "\"hello\" + false",
			expectedError: errors.New("operators must be strings or numbers"),
		},
		{
			input:         "true + false",
			expectedError: errors.New("operators must be strings or numbers"),
		},
		//
		// token.SLASH
		//
		{
			input:    "1 / 2",
			expected: float64(0.5),
		},
		{
			input:    "1 / 0",
			expected: math.Inf(1),
		},
		{
			input:         "1 / false",
			expectedError: errors.New("operands must be numbers"),
		},
		{
			input:         "\"hello\" / false",
			expectedError: errors.New("operands must be numbers"),
		},
		//
		// token.STAR
		//
		{
			input:    "1 * 2",
			expected: float64(2),
		},
		{
			input:    "1 * 2 * 3",
			expected: float64(6),
		},
		{
			input:         "1 * \"hello\"",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
		{
			input:         "\"hello\" * false",
			expected:      nil,
			expectedError: errors.New("operands must be numbers"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			s := scanner.New([]byte(tt.input))
			tokens, errs := s.ScanTokens()
			require.Empty(t, errs)

			p := parser.New(tokens)
			expr, err := p.ParseExpression()
			require.Nil(t, err)

			// expr := stmts[0].(*ast.ExpressionStmt).Expression.(*ast.BinaryExpr)

			var output bytes.Buffer
			i := New(&output)
			result, err := i.VisitBinaryExpr(expr.(*ast.BinaryExpr))
			if tt.expectedError != nil {
				require.Equal(t, tt.expectedError, err)
			} else {
				require.Nil(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestVisitUnaryExpression(t *testing.T) {
	tests := []struct {
		input         string
		expected      any
		expectedError error
	}{
		//
		// token.BANG
		//
		{
			input:    "!true",
			expected: false,
		},
		{
			input:    "!false",
			expected: true,
		},
		//
		// token.MINUS
		//
		{
			input:    "-1",
			expected: float64(-1),
		},
		{
			input:    "-(-1)",
			expected: float64(1),
		},
		{
			input:         "-\"hello\"",
			expected:      nil,
			expectedError: errors.New("invalid cast"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			s := scanner.New([]byte(tt.input))
			tokens, errs := s.ScanTokens()
			require.Empty(t, errs)

			p := parser.New(tokens)
			expr, err := p.ParseExpression()
			require.Nil(t, err)

			var output bytes.Buffer
			i := New(&output)
			result, err := i.VisitUnaryExpr(expr.(*ast.UnaryExpr))
			if tt.expectedError != nil {
				require.Equal(t, tt.expectedError, err)
			} else {
				require.Nil(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
