package logger

import (
	"testing"
)

func TestMaskString(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          `{"username": "john_doe", "password": "123"} <password> 1234 </password> <credentials> 12345 </credentials> base64data: ZnNkZnNkZnNkZnNkZnNkZnNkZnNkZnNmc2RzZGZzZGZkc2ZzZGZzZGZzZGZmc2QK`,
			expectedOutput: `{"username": "john_doe", "password": "***"} <password> **** </password> <credentials> ***** </credentials> base64data: *****`,
		},
		{
			input:          `{"username": "john_doe", "password": "123"} <password> 1234 </password> <credentials> 12345 </credentials> base64data: ZnNkZnNkZnNkZnNkZnNkZnNkZnNkZnNmc2RzZGZzZGZkc2ZzZGZzZGZzZGZmc2QK`,
			expectedOutput: `{"username": "john_doe", "password": "***"} <password> **** </password> <credentials> ***** </credentials> base64data: *****`,
		},
		{
			input:          `{"password": "123"} <password>1234</password> <credentials>12345</credentials>`,
			expectedOutput: `{"password": "***"} <password>****</password> <credentials>*****</credentials>`,
		},
		{
			input:          `<![CDATA[ <password> 123 </password> ]]>, <![CDATA[ <credentials> user:1234 </credentials> ]]>`,
			expectedOutput: `<![CDATA[ <password> *** </password> ]]>, <![CDATA[ <credentials> ********* </credentials> ]]>`,
		},
		{
			input:          `base64data: ABCDEFGH12345==, <password>123</password>, {"password": "1234"}, <credentials>user:1234</credentials>`,
			expectedOutput: `base64data: ABCDEFGH12345==, <password>***</password>, {"password": "****"}, <credentials>*********</credentials>`,
		},
		{
			input:          `"Luong": "12345"`,
			expectedOutput: `"Luong": "*****"`,
		},
		{
			input:          `"   Luong  Nhân viên  ": "12345"`,
			expectedOutput: `"   Luong  Nhân viên  ": "*****"`,
		},
		{
			input:          `" Chi  luong    ": "12345"`,
			expectedOutput: `" Chi  luong    ": "*****"`,
		},
		{
			input:          `"Staff Salary ": "12345"`,
			expectedOutput: `"Staff Salary ": "*****"`,
		},
		{
			input:          `"   Salary    ": "12345"`,
			expectedOutput: `"   Salary    ": "*****"`,
		},
		{
			input:          `"   SALARY    ": "12345"`,
			expectedOutput: `"   SALARY    ": "*****"`,
		},
		{
			input:          `<TT_HoTroKhac2_TheoNgayCong>100000</TT_HoTroKhac2_TheoNgayCong>`,
			expectedOutput: `<TT_HoTroKhac2_TheoNgayCong>******</TT_HoTroKhac2_TheoNgayCong>`,
		},
		{
			input:          `nothing`,
			expectedOutput: `nothing`,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := MaskSensitiveData(tc.input)
			if result != tc.expectedOutput {
				t.Errorf("Expected: %s, Got: %s", tc.expectedOutput, result)
			}
		})
	}
}

func TestMaskWithCustomerData(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          `{"maskedData": "john_doe"}`,
			expectedOutput: `{"maskedData": "********"}`,
		},
	}

	patterns := []string{
		`\"maskedData\"\s*:\s*\"(.*?)\"`,
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := MaskSensitiveData(tc.input, patterns...)
			if result != tc.expectedOutput {
				t.Errorf("Expected: %s, Got: %s", tc.expectedOutput, result)
			}
		})
	}
}
