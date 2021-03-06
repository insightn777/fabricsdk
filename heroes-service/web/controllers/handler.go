package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	//"reflect"
)


func (app *Application)IndexFunc(w http.ResponseWriter, r *http.Request) {

	fmt.Println("check Cookie!!")
	if checkCookie(r) < 0 {
		t, _ := template.ParseFiles("web/templates/index.html")
		t.Execute(w, "test")
	} else {
		http.Redirect(w, r, "/loginSuccess", 302)
	}
}

type loginRequest struct {
	UserID string
	UserPW string
}
func (app *Application)Login(w http.ResponseWriter, r *http.Request) {
	// 클라이언트에서 AJAX사용하여 jsonString 형식으로 입력 ID,PW 넘겨줌

	var loginrequest loginRequest
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &loginrequest)
	inID := loginrequest.UserID
	inPW := loginrequest.UserPW
	fmt.Println("parameter", inID, inPW)

	// login 결과 "fail"로 초기화
	loginreturn := "fail"

	// ask chaincode to get the User Information using inID *************
	// UserInfo :: jsonString from chaincode
	UserInfo, _ := app.Fabric.QueryUserByName(inID)
	fmt.Println("Debug1 : " , UserInfo)


	var userInfo user
	if UserInfo != "[]" {

		// jsonString을 user 구조체 변수 userInfo에다가 unmarshal
		json.Unmarshal([]byte(UserInfo), &userInfo)


		fmt.Println("DEBUG2 : ", userInfo.Name, userInfo.ID, userInfo.Password)
		if userInfo.Password == inPW {
			loginreturn ="success"
		}
	}

	if inID == "admin" && inPW == "admin" {
		loginreturn = "success"
	}

	var cookie http.Cookie
	if loginreturn == "success" {
		/*loginInformation := LoginInfo{
			ID:       inID,
			Password: inPW,
		}*/

		if inID == "admin" {
			cookie = http.Cookie {
				Name : "admin",
				Value : "admin",
				MaxAge: 3600,
				HttpOnly: true,
			}
		} else {

			cookie = http.Cookie{ // cookie 세팅
				//Name:     user_temp.LoginInformation.ID,
				//Value:    user_temp.LoginInformation.Password,
				Name: userInfo.ID,
				Value: userInfo.Password,
				MaxAge:   3600, // 이 시간 초만큼 지나면 쿠키 소멸되는것 확인
				HttpOnly: true,
			}
		}

		http.SetCookie(w, &cookie)
	} //

	w.Write([]byte(loginreturn))
}


func (app *Application)LoginSuccess(w http.ResponseWriter, r *http.Request) {

	user_state := checkCookie(r)	// cookie Name이 "admin" 이면 "1", userID 면 "0"

	fmt.Println("loginSuccess ", user_state)


	// 사용자 이면, 쿠키 Name(사용자 ID)을 가지고 현재 사용자 정보 가져옴
	idInCookie := getCookieName(r)
	// get User information from chaincode using ID ************************
	// var curUser :: jsonString from chaincode
	curUser, _ := app.Fabric.QueryUserByName(idInCookie)

	fmt.Println(curUser)

	if user_state == 1 {
		t, _ := template.ParseFiles("web/templates/admin_menu.html", "web/templates/info.html")
		t.Execute(w, "admin")
	} else if user_state == 0 {
		// 사용자 정보를 jsonString을 페이지 상단에 로그인 정보 표시를 위해 넘겨줌
		t, _ := template.ParseFiles("web/templates/user_menu.html", "web/templates/info.html")
		t.Execute(w, curUser)
	} else { // 사실상 실행될일 없음
		http.Redirect(w, r, "/", 302)
	}
}

//
type UserCount struct {
	Region string
	Count  int
}
func getUserCount(jsonString string) string {	// testInfo 사용 부분
	// 각 동별 주민수 정보 리턴

	var userlists []user
	json.Unmarshal([]byte(jsonString), &userlists)

	var userCountInfo []UserCount

	for i := 0; i < len(userlists); i++ {
		temp_region := userlists[i].Location

		if regionNotExists(temp_region, userCountInfo) {
			// 새로운 지역 필드 추가
			var new_user_count UserCount = UserCount{temp_region, 0}
			userCountInfo = append(userCountInfo, new_user_count)
		}

		// 해당 지역의 주민수 1 증가
		indx := getIndex(temp_region, userCountInfo)
		userCountInfo[indx].Count++
	}

	// userCountInfo 결과 정보를 JsonString으로 리턴
	resultBytes, _ := json.Marshal(userCountInfo)
	resultString := string(resultBytes)

	return resultString
}
func regionNotExists(temp_region string, userCountInfo []UserCount) bool {
	for i := 0; i < len(userCountInfo); i++ {
		if userCountInfo[i].Region == temp_region {
			return false
		}
	}

	return true
}
func getIndex(temp_region string, userCountInfo []UserCount) int {
	for i := 0; i < len(userCountInfo); i++ {
		if userCountInfo[i].Region == temp_region {
			return i
		}
	}

	return -1
}
//


type VoteResult struct {
	Name      string
	StartDate string
	EndDate   string
}
func exportVotesNeededInformation(jsonByteData []byte) string {	// 구조체 사용부분
	// 전체 투표들 정보 중 (user response를 뺀 나머지) 이름, 시작날짜, 종료날짜, articles 정보만 추출하여
	// JSON string으로 반환

	var data []vote
	json.Unmarshal(jsonByteData, &data)

	//var result []VoteResult
	result := make([]VoteResult, len(data))

	for i := 0; i < len(data); i++ {
		var temp VoteResult
		temp.Name = data[i].Votename
		temp.StartDate = data[i].StartDate
		temp.EndDate = data[i].EndDate

		result[i] = temp
	}

	resultBytes, _ := json.Marshal(result)
	resultString := string(resultBytes)

	return resultString
}
func (app *Application)View_result(w http.ResponseWriter, r *http.Request) {
	//handler to show all votes list

	// get All vote results from chaincode *****************
	// votes_result :: chaincode에서 jsonString 으로 받음
	votes_result, _ := app.Fabric.QueryAllVote()

	// 전체 투표들 정보를 받아서, 이정보들 중에서 이름, 날짜, articles 정보만 빼내서 jsonString으로 반환
	resultString := exportVotesNeededInformation([]byte(votes_result))
	fmt.Println("RESULT STRING :: " + resultString)

	// get user count information	***********************
	// usersInfo ::= chaincode에서 jsonString 으로 받음
	usersInfo, _ := app.Fabric.QueryAllUser()
	userCountInfos := getUserCount(usersInfo)
	// 반환값 userCountInfos 는 {지역 이름, 주민 수} 변수들의 배열

	// 주민수와 투표들 정보를 함께 넘겨주기 위해 'splitSeparator' 라는 문자열을 기준으로 합침
	resultString = resultString + "splitSeparator" + userCountInfos

	fmt.Println("VoteReuslt + UserCountInfos : ", resultString)

	t, _ := template.ParseFiles("web/templates/view_result.html")
	t.Execute(w, resultString)
}

func (app *Application)View_vote_result(w http.ResponseWriter, r *http.Request) {
	// handler to show the result of vote you selected

	fmt.Println("voteID: ", r.FormValue("voteID"))
	// voteID를 입력으로 받아와서 체인코드에게 이 ID를 통해서 
	// 해당 투표 결과 정보 가져옴

	// vote_result 변수가 체인코드에서 받아온 특정 투표 결과 jsonString ********************
	// vote_result :: jsonString from chaincode
	vote_result, _ := app.Fabric.QueryVoteByName(r.FormValue("voteID"))


	// 주민 수 정보를 view_result에서 줘야 될지 아니면, view_vote_result에서 줘야 될지..

	// 요건 주민 수 정보를 안 넘길때 사용
	/*t, _ := template.ParseFiles("web/templates/view_vote_result.html")
	t.Execute(w, vote_result)*/

	// 요건 주민 수 정보를 함께 넘길때 사용
	// get user count information	***********************
	// usersInfo ::= chaincode에서 jsonString 으로 받음
	usersInfo, _ := app.Fabric.QueryAllUser()
	userCountInfos := getUserCount(usersInfo)
	// 반환값 userCountInfos 는 {지역 이름, 주민 수} 변수들의 배열

	jsonString := vote_result + "splitSeparator" + userCountInfos

	t, _ := template.ParseFiles("web/templates/view_vote_result.html")
	t.Execute(w, jsonString)
}


// 투표 등록 요청 정보를 저장할 구조체 
type enrollInfo struct {
	VoteName string
	EndDate  string
	Articles []string
	Location string
}
func (app *Application)Enroll_vote(w http.ResponseWriter, r *http.Request) {
	// 투표 등록을 처리하는 핸들러
	// client에서 AJAX로 요청했음

	// 투표 등록 요청 정보를 읽어서 enrollInfo 형식의 구조체로 unmarshal
	var enrollRequest enrollInfo
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &enrollRequest)
	fmt.Println(enrollRequest.VoteName, enrollRequest.EndDate, enrollRequest.Articles, enrollRequest.Location)

	now := strings.Fields(time.Now().String())[0]
	now = strings.Replace(now, "-", "/", -1)

	// enrollRequest정보를 이용하여, CouchDB에 저장된 Vote 구조체 형식으로 생성
	newVote := vote{
		ObjectType : "vote",
		Votename:	enrollRequest.VoteName,
		StartDate: now,
		EndDate:   enrollRequest.EndDate,
		Question:  enrollRequest.Articles,
		//UserResult 생략 -> zero value	// 생략하면 알아서 zero value로 저장됨

		Location: enrollRequest.Location,
	}

	// chaincode에 투표 등록 요청 ******************************
	//vote_result = append(vote_result, newVote)
	txID, err := app.Fabric.InsertVote(newVote.Votename, newVote.StartDate, newVote.EndDate, newVote.Location, newVote.Question[0], newVote.Question[1], newVote.Question[2], newVote.Question[3])
	fmt.Println("received transacton ID : ", txID)
	if(err != nil){
		fmt.Println("failed to insert");
	}else {
		// 클라이언트(웹 페이지)에게 success 문자열 전송
		w.Write([]byte("success!!!!!"))
	}


}

func (app *Application)Membership_manage(w http.ResponseWriter, r *http.Request) {
	// 회원정보 관리 페이지 출력


	// chaincode에서 모든 users 정보 받아옴
	//usersInfo :: jsonString from chaincode
	usersInfo, _ := app.Fabric.QueryAllUser()


	// membership_manage.html 에 이 정보를 넘겨주면, membership_manage.html 파일 javascript code에서 이를 처리하여 출력함
	t, _ := template.ParseFiles("web/templates/membership_manage.html")
	t.Execute(w, usersInfo)
}

func (app *Application)MembershipDelete(w http.ResponseWriter, r *http.Request) {
	// 회원 삭제 기능 핸들러
	// form tag의 submit 이용

	delID := r.FormValue("delID")
	fmt.Println("delete " + delID)

	// ask chaincode to remove delID ******************************
	txID, err := app.Fabric.DeleteUser(delID)
	fmt.Println("received transacton ID : ", txID)
	if(err != nil){
		fmt.Println("failed to delete");
	}else {
		// 특정 user삭제했으면, membership_manage 페이지로 redirect
		http.Redirect(w, r, "/membership_manage", 302)
	}
}

func (app *Application)MembershipModify(w http.ResponseWriter, r *http.Request) {
	// membership_manage 페이지에서 특정 사용자의 '회원정보 수정' 버튼을 눌렀을때

	modifyID := r.FormValue("modifyID")
	fmt.Println("modify name is " + modifyID)

	// modifyName으로 chaincode로 부터 해당 user의 데이터를 가지고 와서
	// userInfo :: jsonString from chaincode 
	userInfo, _ := app.Fabric.QueryUserByName(modifyID)
	fmt.Println("modify" + userInfo)

	// 해당 user 정보들을 가지고 membershipModify.html 템플릿 실행
	t, _ := template.ParseFiles("web/templates/membershipModify.html")
	t.Execute(w, userInfo)
}
func (app *Application)MembershipModifyRequest(w http.ResponseWriter, r *http.Request) {
	// 수정할 정보를 입력하고 '수정하기' 버튼을 눌렀을때 핸들러

	r.ParseForm()

	userName := r.Form["userName"][0]
	userUnum := r.Form["userUnum"][0]
	userID := r.Form["userID"][0]
	newPW := r.Form["newPW"][0]
	newRegion := r.Form["newRegion"][0]

	fmt.Println("DEBUG ", r.Form["userName"][0], r.Form["userUnum"][0], r.Form["userID"][0], r.Form["newPW"][0], r.Form["newRegion"][0])

	// ask chaincode to modify SAME AS Insert
	// **********************************
	txID, err := app.Fabric.InsertUser(userID, userName, newPW, userUnum, newRegion)

	fmt.Println("received transacton ID : ", txID)
	if(err != nil){
		fmt.Println("failed to insert");
	}else {
		// 수정이 완료되면, membership_manage 페이지로 redirect
		http.Redirect(w, r, "/membership_manage", 302)
	}
}

// 계정 찾기 요청 정보를 담을 구조체
type findRequest struct {
	Type string
	ID   string
	Name string
	Unum string
}
func (app *Application)FindAccount(w http.ResponseWriter, r *http.Request) {
	// 클라이언트에서 AJAX로 요청

	var request findRequest
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &request)
	fmt.Println(request.Type, request.ID, request.Name, request.Unum)


	// 일단 모든 users 정보를 가져와 ******************
	// usersInfo :: jsonString from chaincode
	usersInfo, _ := app.Fabric.QueryAllUser()
	var testInfo []user
	json.Unmarshal([]byte(usersInfo), &testInfo)


	result := "-1"

	if request.Type == "id" { // find ID
		fmt.Println("find ID")
		for i := 0; i < len(testInfo); i++ {
			if testInfo[i].Name == request.Name {
				if testInfo[i].SocialNumber == request.Unum {
					result = testInfo[i].ID
				}
			}
		}
	} else { // find PW
		fmt.Println("find PW")
		for i := 0; i < len(testInfo); i++ {
			if testInfo[i].ID == request.ID {
				if testInfo[i].Name == request.Name && testInfo[i].SocialNumber == request.Unum {
					result = testInfo[i].Password
				}
			}
		}

	}

	w.Write([]byte(result))
}



func (app *Application)Signup(w http.ResponseWriter, r *http.Request) { // 회원가입 메소드
	fmt.Println("method : ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/signup.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm() // 넘어온 정보들 파싱
		var obj string
		// 이제 admin의 아이디를 fix 해놓을거기 때문에, id로 사용자가 일반 사용자인지 관리자인지 판단한다.
		// 일단 민준이가 어떻게 값을 줄진 모르겠지만, objectType 을 내 임의로 string으로 넣어봤다.

		fmt.Printf("\n\n\n\n")
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@", r.Body)
		if r.Form["ID"][0] == "admin" {
			obj = "admin"
		} else {
			obj = "user"
		}

		// r.Form[data] 하면 리턴 타입이 []string임. string으로 넣어주기 위해 저렇게 붙여줘야함
		user := user{
			ObjectType: obj,

			ID : r.Form["ID"][0],
			Name: r.Form["Name"][0],
			Password : r.Form["Password"][0],
			SocialNumber: r.Form["IDNumber"][0],
			Location:     r.Form["Region"][0],
			// VoteResult -> zero value
		}

			//회원가입 처리하면 됨
			//입력받은 데이터들을 chaincode로 보내주는 부분 작성하면 됨!
			if user.ObjectType=="user" {
				txID, err := app.Fabric.InsertUser(user.ID, user.Name, user.Password, user.SocialNumber, user.Location)
				//
				if err != nil {
					fmt.Println("Debug : ", "insert fail!");
				} else {
					fmt.Println("DEBUG : ", txID)
				}
			}

		http.Redirect(w, r, "/", 302) // 첫 화면으로 redirect
	}
}


func getCookieName(r *http.Request) string { // cookie 이름 string으로 반환하는 부분.
	h, err := r.Header["Cookie"]
	result := ""
	if !err {
		fmt.Println("there is no cookie!")
	} else {
		//test
		fmt.Println("getCookieName Debug.. r.Header['Cookie'] : ", h)
		for i := 0; i < len(h[0]); i++ {
			if h[0][i] == '=' {
				break
			}
			result += string(h[0][i])
		}
	}
	return result
}

func checkCookie(r *http.Request) int {
	if getCookieName(r) == "admin" {
		return 1
	} else if getCookieName(r) == "" {
		return -1
	} else {
		return 0
	}
}
func (app *Application)Logout(w http.ResponseWriter, r *http.Request) { // 계정 로그아웃 시켜주는 부분. 사실 로그아웃이라기보단 admin이 secret 페이지를 못보게 만드는 부분.

	//쿠키 삭제
	cookie := http.Cookie{
		Name:     getCookieName(r),
		Value:    "",
		MaxAge:   -10,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", 302)
}
///////////////////////////








func getNowDateString() string { // 현재시간을 YYYY/MM/DD 형태의 string으로 가져오기
	now := strings.Fields(time.Now().String())[0]
	//now = strings.Replace(now, "-", "/", -1)
	return now
}


func (app *Application)VotePage(w http.ResponseWriter, r *http.Request) {

	key, ok := r.URL.Query()["Votename"] // GET 방식으로 투표 제목 받기
		if !ok || len(key[0]) < 1 {
			log.Println("Url Param 'Votename' is missing")
			return
		}

	if r.Method == "GET" { // 설문조사를 하는 페이지 띄워주기

		var str string
		var data vote
		str, _ = app.Fabric.QueryVoteByName(key[0])
		json.Unmarshal([]byte(str), &data)
		t, _ := template.ParseFiles("web/templates/votePage.html")
		t.Execute(w, data)

	} else { // 설문조사에 대한 답을 처리하는 메소드.

		r.ParseForm()
		response, _ := strconv.Atoi(r.Form["resultGroup"][0])
		fmt.Printf("\n\n\n\nval: %d, %s\n\n\n", response, r.Form["resultGroup"])
		response++ // 응답한 문항번호를 response에 담음. 만약 여러 문항이면 배열에 담으면 될것.

		err, _ := app.Fabric.InsertVoteResult(key[0], getCookieName(r), strconv.Itoa(response))
		http.Redirect(w, r, "/user_menu", 302)
		fmt.Println("Debug : ", err)
	}
}

func (app *Application)Usermenu(w http.ResponseWriter, r *http.Request) { // 일반 사용자 메뉴
	t, _ := template.ParseFiles("web/templates/user_menu.html")
	t.Execute(w, nil)
}

func (app *Application)UserVoteAll(w http.ResponseWriter, r *http.Request) { // 모든 투표 목록을 보여줌
	var data WebData
	var str, usr string
	var tempVote []vote
	var tempUser user

	usr, _ = app.Fabric.QueryUserByName(getCookieName(r))
	json.Unmarshal([]byte(usr), &tempUser)

	fmt.Println("\n\nLocation : ", tempUser.Location, "\n\n\n")

	str, _ = app.Fabric.QueryVoteByLocation(tempUser.Location)

	fmt.Println("\n\n\nresult : ", str,"\n\n\n")
	json.Unmarshal([]byte(str), &tempVote)

	data.Vote = tempVote
	data.Now = getNowDateString()
	data.UserID = tempUser.ID

	log.Println(data.Vote)

	t, _ := template.ParseFiles("web/templates/userVoteAll.html")
	t.Execute(w, data)
}


func (app *Application)UserVoting(w http.ResponseWriter, r *http.Request) { // 참여하지 않았고 진행중인 투표 목록
	var data WebData
	var str, usr string
	var tempVote []vote
	var tempUser user

	usr, _ = app.Fabric.QueryUserByName(getCookieName(r))
	json.Unmarshal([]byte(usr), &tempUser)

	str, _ = app.Fabric.QueryVoteByLocation(tempUser.Location)
	json.Unmarshal([]byte(str), &tempVote)

	data.Vote = tempVote
	data.Now = getNowDateString()
	data.UserID = tempUser.ID

	t, _ := template.ParseFiles("web/templates/userVoting.html")
	t.Execute(w, data)
}

func (app *Application)UserVoted(w http.ResponseWriter, r *http.Request) { // 투표가 완료됐거나 참여한것만 출력
	var data WebData
	var str, usr string
	var tempVote []vote
	var tempUser user

	usr, _ = app.Fabric.QueryUserByName(getCookieName(r))
        json.Unmarshal([]byte(usr), &tempUser)

        str, _ = app.Fabric.QueryVoteByLocation(tempUser.Location)
        json.Unmarshal([]byte(str), &tempVote)

	data.Vote = tempVote
	data.Now = getNowDateString()
	data.UserID = tempUser.ID

	t, _ := template.ParseFiles("web/templates/userVoted.html")
	t.Execute(w, data)
}

type WebData struct {
	Vote   []vote
	Now    string
	UserID string
}


func (app *Application)CheckID(w http.ResponseWriter, r *http.Request) {
   body, _ := ioutil.ReadAll(r.Body)
   id := string(body)
   result, _ := app.Fabric.QueryUserByName(id)
   fmt.Println("INPUT id : ", id);
   fmt.Println("RESULT : ", result);
   if result == "" {
      w.Write([]byte("true"))
   } else {
      w.Write([]byte("false"))
   }
}

