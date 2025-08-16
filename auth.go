package main

import "net/http"

// 認証ハンドラ
type authHandler struct {
	next http.Handler
}

// 認証に関するメソッド
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// authという名前のクッキーがあるか確認
	// クッキーがない場合、未認証として、loginページへリダイレクト
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// 未認証以外のエラーが発生した場合
		panic(err.Error())
	} else {
		// 成功した場合、ラップされたハンドラを呼び出す
		h.next.ServeHTTP(w, r)
	}
}

// http.Handlerをラップした*authHandlerを生成する、ヘルパー関数
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
