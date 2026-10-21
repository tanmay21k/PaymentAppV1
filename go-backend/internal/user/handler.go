package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

type requestStruct struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		var req requestStruct
		err := json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()
		if err != nil {
			log.Println("err:", err)
			writeMessage(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if !isValidName(req.Firstname) || !isValidName(req.Lastname) {
			log.Println("expected firstname and lastname")
			writeMessage(w, http.StatusBadRequest, "Invalid inputs")
			return
		}
		if !isValidPassword(req.Password) {
			log.Println("password must have more than 6 characters")
			writeMessage(w, http.StatusBadRequest, "Invalid inputs")
			return
		}
		ok := isValidUsername(req.Username)
		if !ok {
			log.Println("username does not match requirements")
			writeMessage(w, http.StatusBadRequest, "Invalid inputs")
			return
		}
	}
}

// writeMessage sends the status and a JSON body like {"msg":"..."}.
// Headers must be set before calling it, because WriteHeader sends them.
func writeMessage(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"msg": msg})
}

func isValidName(name string) bool {
	return name != ""
}

func isValidPassword(password string) bool {
	return utf8.RuneCountInString(password) >= 6
}

// isValidUsername applies the practical email rule: one "@", a local part and
// a domain that each pass their own checks, and at most 254 bytes in total.
func isValidUsername(username string) bool {
	if len(username) > 254 {
		return false
	}
	local, domain, found := strings.Cut(username, "@")
	if !found || strings.Contains(domain, "@") {
		return false
	}
	return isValidLocalPart(local) && isValidDomain(domain)
}

// isValidLocalPart allows 1 to 64 bytes of letters, digits and . _ % + -,
// with no leading, trailing or doubled dot.
func isValidLocalPart(local string) bool {
	if len(local) < 1 || len(local) > 64 {
		return false
	}
	if local[0] == '.' || local[len(local)-1] == '.' {
		return false
	}
	prev := byte(0)
	for i := 0; i < len(local); i++ {
		c := local[i]
		if !isASCIILetter(c) && !isASCIIDigit(c) && !strings.ContainsRune("._%+-", rune(c)) {
			return false
		}
		if c == '.' && prev == '.' {
			return false
		}
		prev = c
	}
	return true
}

// isValidDomain allows 1 to 253 bytes split into labels by dots. Each label is
// letters, digits and hyphens with no leading or trailing hyphen. There must be
// at least two labels, and the last one (the TLD) is at least 2 letters.
func isValidDomain(domain string) bool {
	if len(domain) < 1 || len(domain) > 253 {
		return false
	}
	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return false
	}
	for _, label := range labels {
		if !isValidLabel(label) {
			return false
		}
	}
	tld := labels[len(labels)-1]
	if len(tld) < 2 {
		return false
	}
	for i := 0; i < len(tld); i++ {
		if !isASCIILetter(tld[i]) {
			return false
		}
	}
	return true
}

func isValidLabel(label string) bool {
	if len(label) < 1 || label[0] == '-' || label[len(label)-1] == '-' {
		return false
	}
	for i := 0; i < len(label); i++ {
		c := label[i]
		if !isASCIILetter(c) && !isASCIIDigit(c) && c != '-' {
			return false
		}
	}
	return true
}

func isASCIILetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isASCIIDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
