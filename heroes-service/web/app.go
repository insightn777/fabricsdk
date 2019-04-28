package web

import (
	"fmt"
	"github.com/chainHero/heroes-service/web/controllers"
	"net/http"
	"html/template"
)

/* origin chainHero source
func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/home.html", app.HomeHandler)
	http.HandleFunc("/request.html", app.RequestHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home.html", http.StatusTemporaryRedirect)
	})

	fmt.Println("Listening (http://localhost:8080/) ...")
	http.ListenAndServe(":8080", nil)
}*/

func Serve(app *controllers.Application) {

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.IndexFunc)

	mux.HandleFunc("/signupfrm", func(w http.ResponseWriter, r *http.Request) {
		// 회원가입 기능으로 들어가기 위해 'signup' 버튼 눌렀을때
		t, _ := template.ParseFiles("web/templates/signup.html")
		t.Execute(w, nil)
	})
	mux.HandleFunc("/findaccountfrm", func(w http.ResponseWriter, r *http.Request) {
		// 계정찾기 기능으로 들어가기 위해 'find account' 버튼 눌렀을때
		t, _ := template.ParseFiles("web/templates/find_account.html")
		t.Execute(w, nil)
	})
	mux.HandleFunc("/indexcancle", func(w http.ResponseWriter, r *http.Request) {
		// 어느 페이지에서든 취소를 누르면 main 화면으로 돌림
		http.Redirect(w, r, "/", 302)
	})

	mux.HandleFunc("/funccancle", func(w http.ResponseWriter, r *http.Request) {
		// 로그인 하고 난 뒤, 관리자 기능 페이지에서든, 사용자 기능 페이지에서든지,,
		// 뒤로가기나 취소 버튼을 누르면, loginSuccess action 실행하다도록,
		// 즉, 결과적으로, admin이면 admin menu가 나타나고, user면 user menu가 나타남
		http.Redirect(w, r, "/loginSuccess", 302)
	})

	mux.HandleFunc("/login", app.Login)
	mux.HandleFunc("/loginSuccess", app.LoginSuccess)
	mux.HandleFunc("/signup", app.Signup)
	mux.HandleFunc("/findAccount", app.FindAccount)


	// administrator handler
	mux.HandleFunc("/view_result", app.View_result)	// 전체 투표 결과 목록 조회
	mux.HandleFunc("/view_vote_result", app.View_vote_result) // 특정 투표 결과 조회

	mux.HandleFunc("/onEnrollVote", func(w http.ResponseWriter, r *http.Request) {
		// 투표 등록하기 기능으로 이동
		t, _ := template.ParseFiles("web/templates/enroll_vote.html")
		t.Execute(w, "_")
	})
	mux.HandleFunc("/enroll_vote", app.Enroll_vote)	// 입력한 정보로 투표 등록 기능 실행

	mux.HandleFunc("/membership_manage", app.Membership_manage)
	mux.HandleFunc("/membershipDelete", app.MembershipDelete)
	mux.HandleFunc("/membershipModify", app.MembershipModify)
	mux.HandleFunc("/membershipModifyRequest", app.MembershipModifyRequest)

	mux.HandleFunc("/logout", app.Logout)


	// user 핸들러
	mux.HandleFunc("/usermenu", app.Usermenu)
	mux.HandleFunc("/userVoteAll", app.UserVoteAll)
	mux.HandleFunc("/userVoting", app.UserVoting)
	mux.HandleFunc("/userVoted", app.UserVoted)
	mux.HandleFunc("/votePage", app.VotePage)

	mux.HandleFunc("/checkID", app.CheckID)

	fmt.Println("Listening (http://localhost:8080/) ...")
	http.ListenAndServe(":8080", mux)
	//fmt.Println(http.ListenAndServe(":8080", nil))
}
