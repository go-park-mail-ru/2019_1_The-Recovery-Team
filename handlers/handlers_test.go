package handlers

//type TestCase struct {
//	ID      string
//	Result  *Result
//	IsError bool
//}
//
//type Result struct {
//	Description string
//	SessionID   int
//}

//func GetProfileDummy(w http.ResponseWriter, r *http.Request) {
//	key := r.FormValue("id")
//
//	switch key {
//	case "42":
//		w.WriteHeader(http.StatusOK)
//		_, _ = io.WriteString(w, `{"session_id": 123123123, "description" : "Profile found successfully"}`)
//	case "__not_auth":
//		w.WriteHeader(http.StatusForbidden)
//		_, _ = io.WriteString(w, `{"description" : "Not authorized"}`)
//	case "__not_exist":
//		w.WriteHeader(http.StatusNotFound)
//		_, _ = io.WriteString(w, `{"description" : "Not found"}`)
//	case "__database_error":
//		fallthrough
//	default:
//		w.WriteHeader(http.StatusInternalServerError)
//		_, _ = io.WriteString(w, `{"description" : "Database error"}`)
//	}
//}
//
//func TestGetProfile(t *testing.T) {
//	cases := [] TestCase{
//		TestCase{
//			ID: "42",
//			Result: &Result{
//				SessionID: 123123123,
//				Description: "Profile found successfully",
//			},
//		},
//		TestCase{
//			ID: "__not_auth",
//			Result: &Result{
//				Description: "Not authorized",
//			},
//		},
//		TestCase{
//			ID: "__not_exist",
//			Result: &Result{
//				Description: "Not found",
//			},
//		},
//		TestCase{
//			ID: "__database_error",
//			Result: &Result{
//				Description: "Database error",
//			},
//		},
//	}
//
//	ts := httptest.NewServer(http.HandlerFunc(GetProfileDummy))
//	for caseNum, item := range cases {
//		result, err := c.Checkout(item.ID)
//
//		if err != nil && !item.IsError {
//			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
//		}
//		if err == nil && item.IsError {
//			t.Errorf("[%d] expected error, got nil", caseNum)
//		}
//		if !reflect.DeepEqual(item.Result, result) {
//			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.Result, result)
//		}
//	}
//	ts.Close()
//}
