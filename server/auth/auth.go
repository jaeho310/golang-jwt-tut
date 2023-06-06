package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"time"
)

type UserAndPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MyClaim struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

func Login(c echo.Context) error {
	loginInfo := UserAndPassword{}
	err := c.Bind(&loginInfo)
	if len(loginInfo.Username) < 1 || len(loginInfo.Password) < 1 {
		return echo.ErrBadRequest
	}

	if !VerifyUser(loginInfo) {
		return echo.ErrUnauthorized
	}

	// 본 게시글에 2번 항목에 있는 payload를 채우는 로직입니다.
	// 토큰의 만료시간을 지정해 claims를 지정합니다.
	// 사용자가 로그아웃하는경우 만료시간까지 해당토큰으로 접근하지 못하게 조치해줘야 합니다.
	// go jwt에서 제공하는 standard claims는 아래와 같으며 발행자, 주제 등을 설정할 수 있지만 필수는 아닙니다.
	// 사용자 클레임을 추가할 수도 있으며 예시에서는 username을 추가했습니다.

	//type StandardClaims struct {
	//	Audience  string `json:"aud,omitempty"`
	//	ExpiresAt int64  `json:"exp,omitempty"`
	//	Id        string `json:"jti,omitempty"`
	//	IssuedAt  int64  `json:"iat,omitempty"`
	//	Issuer    string `json:"iss,omitempty"`
	//	NotBefore int64  `json:"nbf,omitempty"`
	//	Subject   string `json:"sub,omitempty"`
	//}

	claims := &MyClaim{
		loginInfo.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
	}

	// 본 게시글의 header에 설명된 암호화 알고리즘을 채우는 부분입니다.
	// 암호화 알고리즘은 HS256, RSA를 많이 사용하며
	// HS256은 대칭키 암호화 방법으로 jwt 토큰 발급시 일반적으로 사용됩니다.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 본 게시글의 signiture 부분입니다.
	// jwt는 secret키를 알고있으면 어느 서버에서나 인증, 인가를 수행할 수 있는 stateless 방식입니다.
	// 즉 디비를 뒤져 세션정보를 찾아 인증하는 형식이 아니므로
	// 여러대의 서버 인스턴스가 secret키만 알고있으면 빠르게 인증가능합니다.
	// secret키는 절대 유출되면 안됩니다.
	// 예시에서는 .env파일을 통해 환경변수를 로드했으나
	// 컨테이너 빌드나, 오케스트레이션 배포시 env를 주입해주거나 신뢰할수 있는 안전한 저장소에 저장합니다.
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": signedToken,
	})
}

func VerifyUser(userInfo UserAndPassword) bool {
	// 예시를 위해 admin만 로그인에 성공시켰지만
	// 회원가입시 encoder 라이브러리를 사용해 암호화된 비밀번호를 db에 저장해놓고
	// 로그인시 입력한 비밀번호가 암호화된 비밀번호가 될수있는지를 확인하는게 일반적입니다.
	// https통신은 비대칭키를 사용한 암호화 방식으로 1차적으로 외부에서 확인할 수 없으나
	// 클라이언트의 요청도 개발자도구에서 보이지 않도록 base64등으로 인코딩하거나 암호화를 하는게 좋습니다.
	if userInfo.Username == "admin" && userInfo.Password == "admin12!@" {
		return true
	}
	return false
}
