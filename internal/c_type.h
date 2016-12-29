/* NH 오픈API 에 포함된 샘플 소스코드에서 복사해서 붙여넣기 된 후,
* GHTS 개발자에 의해서 일부 수정됨.
* 저작권 관련 규정은 원래 샘플 소스코드의 저작권 규정을 따름.
*
* COPIED and FROM 'WmcaIntf.h' in NH OpenAPI sample code and,
* MODIFIED by GHTS Authors.
* LICENSING TERM follows that of original code.
*
* 변수명명규칙
* Go언어와 데이터를 주고 받는 구조체의 멤버 변수는 Go언어와의 호환성을 위해서,
* Go언어에서 public형은 대문자로 시작해야 함.
* C 헤더 파일은 'go tool cgo -godefs'로 바로 Go자료형으로 변환해서 사용하므로,
* 구조체 멤버 필드의 경우 Go언어의 변수명명 규칙을 여기에서도 적용해서 첫 글자를 대문자로 함.
* 그 외 최근 주류 언어인 Java, C#의 관례에 따라 CamelCase를 적용함. */

# include <windef.h>

typedef	BOOL (__stdcall *F_BOOL)();
typedef	BOOL (__stdcall *F_SetServer)(const char* ServerDnsName);
typedef	BOOL (__stdcall *F_SetPort)(const int nPort);
typedef	BOOL (__stdcall *F_SetOption)(const char* Key,const char* Value);
typedef	BOOL (__stdcall *F_Connect)(HWND hWnd,DWORD msg,char mediaType,char userType,const char* pID,const char* pPW,const char* pCertPW);
typedef	BOOL (__stdcall *F_Query)(HWND hWnd,int trId,const char* trCode,const char* pInputData,int inputDataSize,int accountIndex);
typedef	BOOL (__stdcall *F_Attach)(HWND hWnd,const char* pSiseName,const char* pInputCode,int inputCodeSize,int inputCodeTotalSize);
typedef	BOOL (__stdcall *F_Detach)(HWND hWnd,const char* pSiseName,const char* pInputCode,int inputCodeSize,int inputCodeTotalSize);

//----------------------------------------------------------------------//
// WMCA_CONNECTED 로그인 구조체
//----------------------------------------------------------------------//
typedef	struct {
    char 	AccountNo[11];			// 계좌번호
    char	AccountName[40];		// 계좌명
    char	AccountProductCode[3];	// 상품코드
    char	AmnTabCode[4];			// 관리점코드 ?? 도대체 무엇의 약자일까?
    char	ExpirationDate[8];		// 위임만기일
    char	Granted;				// 일괄주문 허용계좌(G:허용)
    char	Filler[189];			// filler ??
} ACCOUNTINFO;

typedef struct {
    char    Date[14];				// 접속시각
    char	ServerName[15];			// 접속서버
    char	UserID[8];				// 접속자ID
    char    AccountCount[3];		// 계좌수
    ACCOUNTINFO	Accountlist	[999];	// 계좌목록
} LOGININFO;

typedef struct {
    int       TrIdNo;
    LOGININFO *LoginInfo;
} LOGINBLOCK;

//----------------------------------------------------------------------//
// WMCA 문자message 구조체
//----------------------------------------------------------------------//
typedef struct  {
    char	MsgCode[5];	// 00000:정상, 기타:비정상(해당 코드값을 이용하여 코딩하지 마세요. 코드값은 언제든지 변경될 수 있습니다.)
    char	UsrMsg[80];
} MSGHEADER;

//----------------------------------------------------------------------//
// WMCA TR 응답 구조체
//----------------------------------------------------------------------//
typedef struct {
    char*	BlockName;
    char*	DataString;
    int	Length;
} RECEIVED;

typedef struct {
    int		  TrIdNo;
    RECEIVED* DataStruct;
} OUTDATABLOCK;

//----------------------------------------------------------------------//
// 주식 현재가 조회 (c1101)
//----------------------------------------------------------------------//

typedef struct { // 기본입력
	char Lang[1];	char _Lang;							// 한영구분
	char Code[6];	char _Code;							// 종목코드
} Tc1101InBlock;

typedef struct { // 종목마스타기본자료
	char Code[6];	char _Code;							// 종목코드
	char Title[13];	char _Title;						// 종목명. 첫자리는 kospi200은 ‘*’, 스타지수종목은 ‘#’. 실제 종목명은 12 byte임
	char MarketPrice[7];	char _MarketPrice;			// 현재가
	char DiffSign[1];	char _DiffSign;					// 등락부호. 0x18 :상한, 0x1E :상승, 0x20 :보함, 0x19 :하한, 0x1F :하락. 등락부호는 시장과 관계없이 동일한 코드체계 사용
	char Diff[6];	char _Diff;							// 등락폭
	char DiffRate[5];	char _DiffRate;					// 등락률
	char OfferPrice[7];	char _OfferPrice;				// 매도 호가
	char BidPrice[7];	char _BidPrice;					// 매수 호가
	char Volume[9];	char _Volume;						// 거래량
	char TrVolRate[6];	char _TrVolRate;				// 거래비율
	char FloatRate[5];	char _FloatRate;				// 유동주회전율
	char TrAmount[9];	char _TrAmount;					// 거래대금
	char UpLmtPrice[7];	char _UpLmtPrice;				// 상한가
	char High[7];	char _High;							// 장중고가
	char Open[7];	char _Open;							// 시가
	char VsOpenSign[1];	char _VsOpenSign;				// 시가대비부호
	char VsOpenDiff[6];	char _VsOpenDiff;				// 시가대비등락폭
	char Low[7];	char _Low;							// 장중저가
	char LowLmtPrice[7];	char _LowLmtPrice;			// 하한가
	char Time[8];	char _Time;				// 호가시간
	char OfferPrice1[7];	char _OfferPrice1;			// 매도 최우선호가
	char OfferPrice2[7];	char _OfferPrice2;			// 매도 차선 호가
	char OfferPrice3[7];	char _OfferPrice3;			// 매도 차차선 호가
	char OfferPrice4[7];	char _OfferPrice4;			// 매도 4차선 호가
	char OfferPrice5[7];	char _OfferPrice5;			// 매도 5차선 호가
	char OfferPrice6[7];	char _OfferPrice6;			// 매도 6차선 호가
	char OfferPrice7[7];	char _OfferPrice7;			// 매도 7차선 호가
	char OfferPrice8[7];	char _OfferPrice8;			// 매도 8차선 호가
	char OfferPrice9[7];	char _OfferPrice9;			// 매도 9차선 호가
	char OfferPrice10[7];	char _OfferPrice10;			// 매도 10차선 호가
	char BidPrice1[7];	char _BidPrice1;				// 매수 최우선 호가
	char BidPrice2[7];	char _BidPrice2;				// 매수 차선 호가
	char BidPrice3[7];	char _BidPrice3;				// 매수 차차선 호가
	char BidPrice4[7];	char _BidPrice4;				// 매수 4차선 호가
	char BidPrice5[7];	char _BidPrice5;				// 매수 5차선 호가
	char BidPrice6[7];	char _BidPrice6;				// 매수 6차선 호가
	char BidPrice7[7];	char _BidPrice7;				// 매수 7차선 호가
	char BidPrice8[7];	char _BidPrice8;				// 매수 8차선 호가
	char BidPrice9[7];	char _BidPrice9;				// 매수 9차선 호가
	char BidPrice10[7];	char _BidPrice10;				// 매수 10차선 호가
	char OfferVolume1[9];	char _OfferVolume1;			// 매도 최우선 잔량
	char OfferVolume2[9];	char _OfferVolume2;			// 매도 차선 잔량
	char OfferVolume3[9];	char _OfferVolume3;			// 매도 차차선 잔량
	char OfferVolume4[9];	char _OfferVolume4;			// 매도 4차선 잔량
	char OfferVolume5[9];	char _OfferVolume5;			// 매도 5차선 잔량
	char OfferVolume6[9];	char _OfferVolume6;			// 매도 6차선 잔량
	char OfferVolume7[9];	char _OfferVolume7;			// 매도 7차선 잔량
	char OfferVolume8[9];	char _OfferVolume8;			// 매도 8차선 잔량
	char OfferVolume9[9];	char _OfferVolume9;			// 매도 9차선 잔량
	char OfferVolume10[9];	char _OfferVolume10;		// 매도 10차선 잔량
	char BidVolume1[9];	char _BidVolume1;				// 매수 최우선 잔량
	char BidVolume2[9];	char _BidVolume2;				// 매수 차선 잔량
	char BidVolume3[9];	char _BidVolume3;				// 매수 차차선 잔량
	char BidVolume4[9];	char _BidVolume4;				// 매수 4차선 잔량
	char BidVolume5[9];	char _BidVolume5;				// 매수 5차선 잔량
	char BidVolume6[9];	char _BidVolume6;				// 매수 6차선 잔량
	char BidVolume7[9];	char _BidVolume7;				// 매수 7차선 잔량
	char BidVolume8[9];	char _BidVolume8;				// 매수 8차선 잔량
	char BidVolume9[9];	char _BidVolume9;				// 매수 9차선 잔량
	char BidVolume10[9];	char _BidVolume10;			// 매수 10차선 잔량
	char OfferVolTot[9];	char _OfferVolTot;			// 총 매도 잔량
	char BidVolTot[9];	char _BidVolTot;				// 총 매수 잔량
	char OfferVolAfterHour[9];	char _OfferVolAfterHour; // 시간외 매도 잔량
	char BidVolAfterHour[9];	char _BidVolAfterHour;	// 시간외 매수 잔량
	char PivotUp2[7];	char _PivotUp2;					// 피봇 2차 저항 : 피봇가 + 전일 고가 – 전일 저가
	char PivotUp1[7];	char _PivotUp1;					// 피봇 1차 저항 : (피봇가 * 2) – 전일 저가
	char PivotPrice[7];	char _PivotPrice;				// 피봇가 : (전일 고가 + 전일 저가 + 전일 종가) / 3
	char PivotDown1[7];	char _PivotDown1;				// 피봇 1차 지지 : (피봇가 * 2) – 전일 고가
	char PivotDown2[7];	char _PivotDown2;				// 피봇 2차 지지 : 피봇가 – 전일고가 + 전일 저가
	char Market[6];	char _Market;						// 코스피/코스닥 구분 : '코스피' , '코스닥'
	char Sector[18];	char _Sector;					// 업종명
	char CapSize[6];	char _CapSize;					// 자본금규모
	char SettleMonth[16];	char _SettleMonth;			// 결산월
	char MarketAction1[16];	char _MarketAction1;		// 시장조치1
	char MarketAction2[16];	char _MarketAction2;		// 시장조치2
	char MarketAction3[16];	char _MarketAction3;		// 시장조치3
	char MarketAction4[16];	char _MarketAction4;		// 시장조치4
	char MarketAction5[16];	char _MarketAction5;		// 시장조치5
	char MarketAction6[16];	char _MarketAction6;		// 시장조치6
	char CircuitBreaker[6];	char _CircuitBreaker;		// 서킷 브레이커 발동 구분
	char NominalPrice[7];	char _NominalPrice;			// 액면가
	char PrevPriceTitle[12];	char _PrevPriceTitle;	// 전일 종가 타이틀 (평가가격, 기준가, 전일종가)
	char PrevPrice[7];	char _PrevPrice;				// 전일종가
	char CollateralValue[7];	char _SubstituteValue;	// 대용가
	char PublicOfferPrice[7];	char _PublicOfferPrice;	// 공모가
	char High5Day[7];	char _High5Day;					// 5일고가
	char Low5Day[7];	char _Low5Day;					// 5일저가
	char High20Day[7];	char _High20Day;				// 20일고가
	char Low20Day[7];	char _Low20Day;					// 20일저가
	char High1Year[7];	char _High1Year;				// 52주최고가
	char High1YearDate[4];	char _High1YearDate;		// 52주최고가일
	char Low1Year[7];	char _Low1Year;					// 52주최저가
	char Low1YearDate[4];	char _Low1YearDate;			// 52주최저가일
	char FloatVolume[8];	char _FloatVolume;			// 유동주식수
	char ListVolBy1000[12];	char _ListVolBy1000;		// 상장주식수. 1000주 단위?
	char MarketCapital[9];	char _MarketCapital;		// 시가총액
	char TraderInfoTime[5];	char _TraderInfoTime;		// 거래원 정보 최종 수신 시간
	char Seller1[6];	char _Seller1;					// 매도 거래원1
	char Buyer1[6];	char _Buyer1;						// 매수 거래원1
	char Seller1Volume[9];	char _Seller1Volume;		// 매도 거래량1
	char Buyer1Volume[9];	char _Buyer1Volume;			// 매수 거래량1
	char Seller2[6];	char _Seller2;					// 매도 거래원2
	char Buyer2[6];	char _Buyer2;						// 매수 거래원2
	char Seller2Volume[9];	char _Seller2Volume;		// 매도 거래량2
	char Buyer2Volume[9];	char _Buyer2Volume;			// 매수 거래량2
	char Seller3[6];	char _Seller3;					// 매도 거래원3
	char Buyer3[6];	char _Buyer3;						// 매수 거래원3
	char Seller3Volume[9];	char _Seller3Volume;		// 매도 거래량3
	char Buyer3Volume[9];	char _Buyer3Volume;			// 매수 거래량3
	char Seller4[6];	char _Seller4;					// 매도 거래원4
	char Buyer4[6];	char _Buyer4;						// 매수 거래원4
	char Seller4Volume[9];	char _Seller4Volume;		// 매도 거래량4
	char Buyer4Volume[9];	char _Buyer4Volume;			// 매수 거래량4
	char Seller5[6];	char _Seller5;					// 매도 거래원5
	char Buyer5[6];	char _Buyer5;						// 매수 거래원5
	char Seller5Volume[9];	char _Seller5Volume;		// 매도 거래량5
	char Buyer5Volume[9];	char _Buyer5Volume;			// 매수 거래량5
	char ForeignSellVolume[9];	char _ForeignSellVolume; // 외국인 매도 거래량
	char ForeignBuyVolume[9];	char _ForeignBuyVolume;	// 외국인 매수 거래량
	char ForeignTime[6];	char _ForeignTime;			// 외국인 시간 ???
	char ForeignHoldingRate[5];	char _ForeignHoldingRate; // 외국인 지분율
	char SettleDate[4];	char _SettleDate;				// 결제일
	char DebtPercent[5];	char _DebtPercent;			// 잔고 비율(%)
	char RightsIssueDate[4];	char _RightsIssueDate;	// 유상 기준일
	char BonusIssueDate[4];	char _BonusIssueDate;		// 무상 기준일
	char RightsIssueRate[5];	char _RightsIssueRate;	// 유상 배정비율
	char BonusIssueRate[5];	char _BonusIssueRate;		// 무상 배정비율
	char ForeignFloatVol[10];	char _ForeignFloatVol;	// 외국인 변동주 수
	char TreasuryStock[1];	char _TreasuryStock;		// 당일 자사주 신청 여부  1: 자사주 신청
	char IpoDate[8];	char _IpoDate;					// 상장일
	char MajorHoldRate[5];	char _MajorHoldRate;		// 대주주지분율
	char MajorHoldInfoDate[6];	char _MajorHoldInfoDate; // 대주주지분일자
	char FourLeafClover[1];	char _FourLeafClover;		// 네잎클로버 종목 여부 1: 네잎클로버 종목
	char MarginRate[1];	char _MarginRate;				// 증거금율
	char Capital[9];	char _Capital;					// 자본금
	char SellTotalSum[9];	char _SellTotalSum;			// 전체 거래원 매도 합계
	char BuyTotalSum[9];	char _BuyTotalSum;			// 전체 거래원 매수 합계
	char Title2[21];	char _Title2;					// 종목명2. 앞에 한자리를 제외하고 18byte가 종목명
	char BackdoorListing[1];	char _BackdoorListing;	// 우회상장여부
	char FloatRate2[6];	char _FloatRate2;				// 유동주회전율2
	char Market2[6];	char _Market2;					// 코스피 구분 ?? 앞에 나왔는 데...
	char DebtTrDate[4];	char _DebtTrDate;				// 공여율기준일
	char DebtTrPercent[5];	char _DebtTrPercent;		// 공여율(%)
	char PER[5];	char _PER;							// PER
	char DebtLimit[1];	char _DebtLimit;				// 종목별신용한도
	char WeightAvgPrice[7];	char _WeightAvgPrice;		// 가중가
	char ListedVolume[12];	char _ListedVolume;			// 상장주식 수  _주
	char AddListing[12];	char _AddListing;			// 추가상장 주식 수
	char Comment[100];	char _Comment;					// 종목 comment
	char PrevVolume[9];	char _PrevVolume;				// 전일 거래량
	char VsPrevSign[1];	char _VsPrevSign;				// 전일대비 등락부호
	char VsPrevDiff[6];	char _VsPrevDiff;				// 전일대비 등락폭
	char High1Year2[7];	char _High1Year2;				// 연중 최고가 (52주 최고가와 중복 아닌가?
	char High1YearDate2[4];	char _High1YearDate2;		// 연중 최고가일
	char Low1Year2[7];	char _Low1Year2;				// 연중 최저가
	char Low1YearDate2[4];	char _Low1YearDate2;		// 연중 최저가일
	char ForeignHoldQty[15];	char _ForeignHoldQty;	// 외국인 보유 주식수
	char ForeignLmtPercent[5];	char _ForeignLmtPercent; // 외국인 한도율(%)
	char TrUnitVolume[5];	char _TrUnitVolume;			// 매매 수량 단위
	char DarkPoolOfferBid[1];	char _DarkPoolOfferBid; // 경쟁대량방향구분. 0: 해당없음, 1: 매도, 2: 매수
	char DarkPoolExist[1];	char _DarkPoolExist;		// 대량매매구분. 1: 대량매매有, 0:대량매매無
} Tc1101OutBlock;

typedef struct { // 변동거래량자료,[반복]
	char Time[8];	char _Time;							// 시간
	char MarketPrice[7];	char _MarketPrice;			// 현재가
	char DiffSign[1];	char _DiffSign;					// 등락부호
	char Diff[6];	char _Diff;							// 등락폭
	char OfferPrice[7];	char _OfferPrice;				// 매도 호가
	char BidPrice[7];	char _BidPrice;					// 매수 호가
	char DiffVolume[8];	char _DiffVolume;				// 변동거래량
	char Volume[9];	char _Volume;						// 거래량
} Tc1101OutBlock2;

typedef struct { // 종목지표
	char SyncOfferBid[1];	char _SyncOfferBid;			// 동시호가 구분.  0:동시호가 아님, 1:동시호가, 2:동시호가연장, 3:시가범위연장, 4:종가범위연장, 5:배분개시, 6:변동성 완화장치 발동
	char EstmPrice[7];	char _EstmPrice;				// 예상체결가
	char EstmSign[1];	char _EstmSign;					// 예상체결 부호
	char EstmDiff[6];	char _EstmDiff;					// 예상체결 등락폭
	char EstmDiffRate[5];	char _EstmDiffRate;			// 예상체결 등락률
	char EstmVolume[9];	char _EstmVol;					// 예상체결수량
	char ECN_InfoExist[1];	char _ECN_InfoExist;		// ECN정보 유무 구분 (우리나라에는 ECN이 아직 없을텐데...)
	char ECN_PrevPrice[9];	char _ECN_PrevPrice;		// ECN 전일종가
	char ECN_DiffSign[1];	char _ECN_DiffSign;			// ECN 부호
	char ECN_Diff[9];	char _ECN_Diff;					// ECN 등락폭
	char ECN_DiffRate[5];	char _ECN_DiffRate;			// ECN 등락률
	char ECN_Volume[10];	char _ECN_Volume;			// ECN 체결수량
	char VsECN_EstmSign[1];	char _VsECN_EstmSign; 		// ECN대비 예상 체결 부호
	char VsECN_EstmDiff[6];	char _VsECN_EstmDiff;		// ECN대비 예상 체결 등락폭
	char VsECN_EstmDiffRate[5];	char _ECN_EstmDiffRate;	// ECN대비 예상 체결 등락률
} Tc1101OutBlock3;

//----------------------------------------------------------------------//
// ETF 현재가 조회 (c1151)
//----------------------------------------------------------------------//
typedef struct { // 기본입력
	char Lang[1];	char _Lang;							// 한영구분. 기본값 'K'
	char Code[6];	char _Code;							// 종목코드
} Tc1151InBlock;

typedef struct { // 종목마스타기본자료
	char Code[6];	char _Code;							// 종목코드
	char Title[13];	char _Title;						// 종목명
	char MarketPrice[7];	char _MarketPrice;			// 현재가
	char DiffSign[1];	char _DiffSign;					// 등락부호
	char Diff[6];	char _Diff;							// 등락폭
	char DiffRate[5];	char _DiffRate;					// 등락률
	char OfferPrice[7];	char _OfferPrice;				// 매도 호가
	char BidPrice[7];	char _BidPrice;					// 매수 호가
	char Volume[9];	char _Volume;						// 거래량
	char TrVolRate[6];	char _TrVolRate;				// 거래비율
	char FloatVolRate[5];	char _FloatVolRate;			// 유동주회전율
	char TrAmount[9];	char _TrAmount;					// 거래대금
	char UpLmtPrice[7];	char _UpLmtPrice;				// 상한가
	char High[7];	char _High;							// 장중고가
	char Open[7];	char _Open;							// 시가
	char VsOpenSign[1];	char _VsOpenSign;				// 시가대비부호
	char VsOpenDiff[6];	char _VsOpenDiff;				// 시가대비등락폭
	char Low[7];	char _Low;							// 장중저가
	char LowLmtPrice[7];	char _LowLmtPrice;			// 하한가
	char Time[8];	char _Time;							// 호가시간
	char OfferPrice1[7];	char _OfferPrice1;			// 매도 최우선 호가
	char OfferPrice2[7];	char _OfferPrice2;			// 매도 차선 호가
	char OfferPrice3[7];	char _OfferPrice3;			// 매도 차차선 호가
	char OfferPrice4[7];	char _OfferPrice4;			// 매도 4차선 호가
	char OfferPrice5[7];	char _OfferPrice5;			// 매도 5차선 호가
	char OfferPrice6[7];	char _OfferPrice6;			// 매도 6차선 호가
	char OfferPrice7[7];	char _OfferPrice7;			// 매도 7차선 호가
	char OfferPrice8[7];	char _OfferPrice8;			// 매도 8차선 호가
	char OfferPrice9[7];	char _OfferPrice9;			// 매도 9차선 호가
	char OfferPrice10[7];	char _OfferPrice10;			// 매도 10차선 호가
	char BidPrice1[7];	char _BidPrice1;				// 매수 최우선 호가
	char BidPrice2[7];	char _BidPrice2;				// 매수 차선 호가
	char BidPrice3[7];	char _BidPrice3;				// 매수 차차선 호가
	char BidPrice4[7];	char _BidPrice4;				// 매수 4차선 호가
	char BidPrice5[7];	char _BidPrice5;				// 매수 5차선 호가
	char BidPrice6[7];	char _BidPrice6;				// 매수 6차선 호가
	char BidPrice7[7];	char _BidPrice7;				// 매수 7차선 호가
	char BidPrice8[7];	char _BidPrice8;				// 매수 8차선 호가
	char BidPrice9[7];	char _BidPrice9;				// 매수 9차선 호가
	char BidPrice10[7];	char _BidPrice10;				// 매수 10차선 호가
	char OfferVolume1[9];	char _OfferVolume1;			// 매도 최우선 잔량
	char OfferVolume2[9];	char _OfferVolume2;			// 매도 차선 잔량
	char OfferVolume3[9];	char _OfferVolume3;			// 매도 차차선 잔량
	char OfferVolume4[9];	char _OfferVolume4;			// 매도 4차선 잔량
	char OfferVolume5[9];	char _OfferVolume5;			// 매도 5차선 잔량
	char OfferVolume6[9];	char _OfferVolume6;			// 매도 6차선 잔량
	char OfferVolume7[9];	char _OfferVolume7;			// 매도 7차선 잔량
	char OfferVolume8[9];	char _OfferVolume8;			// 매도 8차선 잔량
	char OfferVolume9[9];	char _OfferVolume9;			// 매도 9차선 잔량
	char OfferVolume10[9];	char _OfferVolume10;		// 매도 10차선 잔량
	char BidVolume1[9];	char _BidVolume1;				// 매수 최우선 잔량
	char BidVolume2[9];	char _BidVolume2;				// 매수 차선 잔량
	char BidVolume3[9];	char _BidVolume3;				// 매수 차차선 잔량
	char BidVolume4[9];	char _BidVolume4;				// 매수 4차선 잔량
	char BidVolume5[9];	char _BidVolume5;				// 매수 5차선 잔량
	char BidVolume6[9];	char _BidVolume6;				// 매수 6차선 잔량
	char BidVolume7[9];	char _BidVolume7;				// 매수 7차선 잔량
	char BidVolume8[9];	char _BidVolume8;				// 매수 8차선 잔량
	char BidVolume9[9];	char _BidVolume9;				// 매수 9차선 잔량
	char BidVolume10[9];	char _BidVolume10;			// 매수 10차선 잔량
	char OfferVolTot[9];	char _OfferVolTot;			// 총 매도 잔량
	char BidVolTot[9];	char _BidVolTot;				// 총 매수 잔량
	char OfferVolAfterHour[9];	char _OfferVolAfterHour; // 시간외 매도 잔량
	char BidVolAfterHour[9];	char _BidVolAfterHour;	// 시간외 매수 잔량
	char PivotUp2[7];	char _PivotUp2;					// 피봇 2차 저항
	char PivotUp1[7];	char _PivotUp1;					// 피봇 1차 저항
	char PivotPrice[7];	char _PivotPrice;				// 피봇가
	char PivotDown1[7];	char _PivotDown1;				// 피봇 1차 지지
	char PivotDown2[7];	char _PivotDown2;				// 피봇 2차 지지
	char Market[6];	char _Market;						// 코스피/코스닥 구분
	char Sector[18];	char _Sector;					// 업종명
	char CapSize[6];	char _CapSize;					// 자본금규모
	char SettleMonth[16];	char _SettleMonth;			// 결산월
	char MarketAction1[16];	char _MarketAction1;		// 시장조치1
	char MarketAction2[16];	char _MarketAction2;		// 시장조치2
	char MarketAction3[16];	char _MarketAction3;		// 시장조치3
	char MarketAction4[16];	char _MarketAction4;		// 시장조치4
	char MarketAction5[16];	char _MarketAction5;		// 시장조치5
	char MarketAction6[16];	char _MarketAction6;		// 시장조치6
	char CircuitBreaker[6];	char _CircuitBreaker;		// 서킷 브레이커 구분
	char NominalPrice[7];	char _NominalPrice;			// 액면가
	char PrevPriceTitle[12];	char _PrevPriceTitle;	// 전일 종가 타이틀
	char PrevPrice[7];	char _PrevPrice;				// 전일종가
	char MortgageValue[7];	char _MortgageValue;		// 대용가
	char PublicOfferPrice[7];	char _PublicOfferPrice;	// 공모가
	char High5Day[7];	char _High5Day;					// 5일고가
	char Low5Day[7];	char _Low5Day;					// 5일저가
	char High20Day[7];	char _High20Day;				// 20일고가
	char Low20Day[7];	char _Low20Day;					// 20일저가
	char High1Year[7];	char _High1Year;				// 52주최고가
	char High1YearDate[4];	char _High1YearDate;		// 52주최고가일
	char Low1Year[7];	char _Low1Year;					// 52주최저가
	char Low1YearDate[4];	char _Low1YearDate;			// 52주최저가일
	char FloatVolume[8];	char _FloatVolume;			// 유동주식수
	char ListVolBy1000[12];	char _ListVolBy1000;		// 상장주식수_천주
	char MarketCapital[9];	char _MarketCapital;		// 시가총액
	char TraderInfoTime[5];	char _TraderInfoTime;		// 거래원 정보 최종 수신 시간
	char Seller1[6];	char _Seller1;					// 매도 거래원1
	char Buyer1[6];	char _Buyer1;						// 매수 거래원1
	char Seller1Volume[9];	char _Seller1Volume;		// 매도 거래량1
	char Buyer1Volume[9];	char _Buyer1Volume;			// 매수 거래량1
	char Seller2[6];	char _Seller2;					// 매도 거래원2
	char Buyer2[6];	char _Buyer2;						// 매수 거래원2
	char Seller2Volume[9];	char _Seller2Volume;		// 매도 거래량2
	char Buyer2Volume[9];	char _Buyer2Volume;			// 매수 거래량2
	char Seller3[6];	char _Seller3;					// 매도 거래원3
	char Buyer3[6];	char _Buyer3;						// 매수 거래원3
	char Seller3Volume[9];	char _Seller3Volume;		// 매도 거래량3
	char Buyer3Volume[9];	char _Buyer3Volume;			// 매수 거래량3
	char Seller4[6];	char _Seller4;					// 매도 거래원4
	char Buyer4[6];	char _Buyer4;						// 매수 거래원4
	char Seller4Volume[9];	char _Seller4Volume;		// 매도 거래량4
	char Buyer4Volume[9];	char _Buyer4Volume;			// 매수 거래량4
	char Seller5[6];	char _Seller5;					// 매도 거래원5
	char Buyer5[6];	char _Buyer5;						// 매수 거래원5
	char Seller5Volume[9];	char _Seller5Volume;		// 매도 거래량5
	char Buyer5Volume[9];	char _Buyer5Volume;			// 매수 거래량5
	char ForeignSellVolume[9];	char _ForeignSellVolume; // 외국인 매도 거래량
	char ForeignBuyVolume[9];	char _ForeignBuyVolume;	// 외국인 매수 거래량
	char ForeignTime[6];	char _ForeignTime;			// 외국인 시간 ???
	char ForeignHoldingRate[5];	char _ForeignHoldingRate; // 외국인 지분율
	char SettleDate[4];	char _SettleDate;				// 결제일
	char DebtPercent[5];	char _DebtPercent;			// 잔고비율(%)
	char RightsIssueDate[4];	char _RightsIssueDate;	// 유상기준일
	char BonusIssueDate[4];	char _BonusIssueDate;		// 무상기준일
	char RightsIssueRate[5];	char _RightsIssueRate;	// 유상배정비율
	char BonusIssueRate[5];	char _BonusIssueRate;		// 무상배정비율
	char IpoDate[8];	char _IpoDate;					// 상장일
	char ListedVolume[12];	char _ListedVolume;			// 상장주식수_주
	char SellTotalSum[9];	char _SellTotalSum;			// 전체 거래원 매도 합계
	char BuyTotalSum[9];	char _BuyTotalSum;			// 전체 거래원 매수 합계
} Tc1151OutBlock;

typedef struct { // 변동거래량자료
	char Time[8];	char _Time;							// 시간
	char MarketPrice[7];	char _MarketPrice;			// 현재가
	char DiffSign[1];	char _DiffSign;					// 등락부호
	char Diff[6];	char _Diff;							// 등락폭
	char OfferPrice[7];	char _OfferPrice;				// 매도 호가
	char BidPrice[7];	char _BidPrice;					// 매수 호가
	char DiffVolume[8];	char _DiffVolume;				// 변동거래량
	char Volume[9];	char _Volume;						// 거래량
} Tc1151OutBlock2;

typedef struct { // 예상체결
	char SyncOfferBid[1];	char _SyncOfferBid;			// 동시 호가 구분
	char EstmPrice[7];	char _EstmPrice;				// 예상 체결가
	char EstmSign[1];	char _EstmSign;					// 예상 체결 부호
	char EstmDiff[6];	char _EstmDiff;					// 예상 체결 등락폭
	char EstmDiffRate[5];	char _EstmDiffRate;			// 예상 체결 등락률
	char EstmVolume[9];	char _EstmVolume;				// 예상체결 수량
} Tc1151OutBlock3;

typedef struct { // ETF자료
	char ETF[1];	char _ETF;							// ETF 구분
	char NAV[9];	char _NAV;							// 장중/최종 NAV
	char DiffSign[1];	char _DiffSign;					// NAV 등락 부호
	char Diff[9];	char _Diff;							// NAV 등락폭
	char PrevNAV[9];	char _PrevNAV;					// 전일 NAV
	char DivergeRate[9];	char _DivergeRate;			// 괴리율
	char DivergeSign[1];	char _DivergeSign;			// 괴리율 부호
	char DividendPerCU[18];	char _DividendPerCU;		// CU(Creation Unit : 설정단위)당 현금 배당액(원)
	char ConstituentNo[4];	char _ConstituentNo;		// 구성 종목수
	char NAVBy100Million[7];	char _NAVBy100Million;	// 순자산총액(억원)
	char TrackingErrRate[9];	char _TrackingErrRate;	// 추적오차율
	char LP_OfferVolume1[9];	char _LP_OfferVolume1;	// LP 매도 최우선 잔량
	char LP_OfferVolume2[9];	char _LP_OfferVolume2;	// LP 매도 차선 잔량
	char LP_OfferVolume3[9];	char _LP_OfferVolume3;	// LP 매도 차차선 잔량
	char LP_OfferVolume4[9];	char _LP_OfferVolume4;	// LP 매도 4차선 잔량
	char LP_OfferVolume5[9];	char _LP_OfferVolume5;	// LP 매도 5차선 잔량
	char LP_OfferVolume6[9];	char _LP_OfferVolume6;	// LP 매도 6차선 잔량
	char LP_OfferVolume7[9];	char _LP_OfferVolume7;	// LP 매도 7차선 잔량
	char LP_OfferVolume8[9];	char _LP_OfferVolume8;	// LP 매도 8차선 잔량
	char LP_OfferVolume9[9];	char _LP_OfferVolume9;	// LP 매도 9차선 잔량
	char LP_OfferVolume10[9];	char _LP_OfferVolume10;	// LP 매도 10차선 잔량
	char LP_BidVolume1[9];	char _LP_BidVolume1;		// LP 매수 최우선 잔량
	char LP_BidVolume2[9];	char _LP_BidVolume2;		// LP 매수 차선 잔량
	char LP_BidVolume3[9];	char _LP_BidVolume3;		// LP 매수 차차선 잔량
	char LP_BidVolume4[9];	char _LP_BidVolume4;		// LP 매수 4차선 잔량
	char LP_BidVolume5[9];	char _LP_BidVolume5;		// LP 매수 5차선 잔량
	char LP_BidVolume6[9];	char _LP_BidVolume6;		// LP 매수 6차선 잔량
	char LP_BidVolume7[9];	char _LP_BidVolume7;		// LP 매수 7차선 잔량
	char LP_BidVolume8[9];	char _LP_BidVolume8;		// LP 매수 8차선 잔량
	char LP_BidVolume9[9];	char _LP_BidVolume9;		// LP 매수 9차선 잔량
	char LP_BidVolume10[9];	char _LP_BidVolume10;		// LP 매수 10차선 잔량
	char TrackingMethod[8];	char _TrackingMethod;		// ETF 복제 방법 구분 코드
	char ETF_Type[6];	char _ETF_Type;					// ETF 상품 유형 코드
} Tc1151OutBlock4;

typedef struct { // 베이스지수자료
	char SectorCode[2];	char _SectorCode;				// 업종코드
	char IndexCode[4];	char _IndexCode;				// 지수코드
	char IndexName[20];	char _IndexName;				// 지수명
	char KP200Index[8];	char _KP200Index;				// 지수
	char KP200Sign[1];	char _KP200Sign;				// 등락부호
	char KP200Diff[8];	char _KP200Diff;				// 등락폭
	char BondIndex[10];	char _BondIndex;				// 채권지수
	char BondSign[1];	char _BondSign;					// 채권등락부호
	char BondDiff[10];	char _BondDiff;					// 채권등락폭
	char ForeignIndexSymbol[12];	char _ForeignIndexSymbol; // 해외지수심볼
	char EtcSectorCode[3];	char _EtcSectorCode;		// 기타업종코드
	char BondIndexCode[6];	char _BondIndexCode;		// 채권지수코드
	char BondDetailCode[1];	char _BondDetailCode;		// 채권지수세부코드
} Tc1151OutBlock5;

//----------------------------------------------------------------------//
// 실시간 정보 질의 : 종목 코드 (h1/k3, h2/k4, h3/k5, j8/k8, j1/j0)
//----------------------------------------------------------------------//
typedef struct { // 입력
	char Code[6];				// 종목코드
} TCodeInBlock;

//----------------------------------------------------------------------//
// 실시간 정보 질의 : 업종 코드 (u1/k1)
//----------------------------------------------------------------------//
typedef struct { // 입력
	char SectorCode[2];	char _SectorCode;			// 업종코드
} TSectorCodeInBlock;

//----------------------------------------------------------------------//
// 코스피/코스닥 호가 잔량 (h1/k3)
//----------------------------------------------------------------------//
typedef struct { // 출력
	char skip[3];				// 앞쪽 3바이트는 패킷유형과 압축구분이므로 skip
	char Code[6];				// 종목코드
	char Time[8];				// 시간
	char OfferPrice1[7];		// 매도 호가
	char BidPrice1[7];			// 매수 호가
	char OfferVolume1[9];		// 매도 호가잔량
	char BidVolume1[9];			// 매수 호가잔량
	char OfferPrice2[7];		// 차선 매도 호가
	char BidPrice2[7];			// 차선 매수 호가
	char OfferVolume2[9];		// 차선 매도 호가잔량
	char BidVolume2[9];			// 차선 매수 호가잔량
	char OfferPrice3[7];		// 차차선 매도 호가
	char BidPrice3[7];			// 차차선 매수 호가
	char OfferVolume3[9];		// 차차선 매도 호가잔량
	char BidVolume3[9];			// 차차선 매수 호가잔량
	char OfferPrice4[7];		// 4차선 매도 호가
	char BidPrice4[7];			// 4차선 매수 호가
	char OfferVolume4[9];		// 4차선 매도 호가잔량
	char BidVolume4[9];			// 4차선 매수 호가잔량
	char OfferPrice5[7];		// 5차선 매도 호가
	char BidPrice5[7];			// 5차선 매수 호가
	char OfferVolume5[9];		// 5차선 매도 호가잔량
	char BidVolume5[9];			// 5차선 매수 호가잔량
	char OfferVolumeTotal[9];	// 총매도호가잔량
	char BidVolumeTotal[9];		// 총매수호가잔량
	char OfferPrice6[7];		// 6차선 매도 호가
	char BidPrice6[7];			// 6차선 매수 호가
	char OfferVolume6[9];		// 6차선 매도 호가잔량
	char BidVolume6[9];			// 6차선 매수 호가잔량
	char OfferPrice7[7];		// 7차선 매도 호가
	char BidPrice7[7];			// 7차선 매수 호가
	char OfferVolume7[9];		// 7차선 매도 호가잔량
	char BidVolume7[9];			// 7차선 매수 호가잔량
	char OfferPrice8[7];		// 8차선 매도 호가
	char BidPrice8[7];			// 8차선 매수 호가
	char OfferVolume8[9];		// 8차선 매도 호가잔량
	char BidVolume8[9];			// 8차선 매수 호가잔량
	char OfferPrice9[7];		// 9차선 매도 호가
	char BidPrice9[7];			// 9차선 매수 호가
	char OfferVolume9[9];		// 9차선 매도 호가잔량
	char BidVolume9[9];			// 9차선 매수 호가잔량
	char OfferPrice10[7];		// 10차선 매도 호가
	char BidPrice10[7];			// 10차선 매수 호가
	char OfferVolume10[9];		// 10차선 매도 호가잔량
	char BidVolume10[9];		// 10차선 매수 호가잔량
	char Volume[9];				// 누적거래량
} Th1k3OutBlock;

//----------------------------------------------------------------------//
// 코스피/코스닥 시간외 호가 잔량 (h2/k4)
//----------------------------------------------------------------------//
typedef struct { // 출력
	char skip[3];				// 앞쪽 3바이트는 패킷유형과 압축구분이므로 skip
	char Code[6];				// 종목코드
	char Time[8];				// 시간
	char OfferVolume[9];		// 총 매도 호가잔량
	char BidVolume[9];			// 총 매수 호가잔량
} Th2k4OutBlock;

//----------------------------------------------------------------------//
// 코스피/코스닥 예상 호가 잔량 (h3/k5)
//----------------------------------------------------------------------//
typedef struct { // 출력
	char skip[3];				// 앞쪽 3바이트는 패킷유형과 압축구분이므로 skip
	char Code[6];				// 종목코드
	char Time[8];				// 시간
	char SyncOfferBid[1];		// 동시구분
	char EstmPrice[7];			// 예상체결가
	char EstmDiffSign[1];		// 예상등락부호
	char EstmDiff[6];			// 예상등락폭
	char EstmDiffRate[5];		// 예상등락률
	char EstmVolume[9];			// 예상체결수량
	char OfferPrice[7];			// 매도 호가
	char BidPrice[7];			// 매수 호가
	char OfferVolume[9];		// 매도 호가잔량
	char BidVolume[9];			// 매수 호가잔량
} Th3k5OutBlock;

//----------------------------------------------------------------------//
// 코스피/코스닥 체결 (j8/k8)
//----------------------------------------------------------------------//
typedef struct { // 출력
	char skip[3];									// 앞쪽 3바이트는 패킷유형과 압축구분이므로 skip
	char Code[6];	char _Code;						// 종목코드
	char Time[8];	char _Time;						// 시간
	char DiffSign[1];	char _DiffSign;				// 등락부호
	char Diff[6];	char _Diff;						// 등락폭
	char MarketPrice[7];	char _MarketPrice;		// 현재가
	char DiffRate[5];	char _DiffRate;				// 등락률
	char High[7];	char _High;						// 고가
	char Low[7];	char _Low;						// 저가
	char OfferPrice[7];	char _OfferPrice;			// 매도 호가
	char BidPrice[7];	char _BidPrice;				// 매수 호가
	char Volume[9];	char _Volume;					// 거래량
	char VsPrevVolRate[6];	char _VsPrevVolRate;	// 거래량전일비
	char DiffVolume[8];	char _DiffVolume;			// 변동거래량
	char TrAmount[9];	char _TrAmount;				// 거래대금
	char Open[7];	char _Open;						// 시가
	char WeightAvgPrice[7];	char _WeightAvgPrice;	// 가중평균가
	char Market[1];	char _Market;					// 장구분 ('0': 코스피, '1': 코스닥)
} Tj8OutBlock;

typedef struct {    // 출력
    char skip[3];									// 앞쪽 3바이트는 패킷유형과 압축구분이므로 skip
	char Code[6];	char _Code;                     //종목코드
	char Time[8];	char _Time;                     //시간
	char MarketPrice[7];	char _MarketPrice;      //현재가
	char DiffSign[1];	char _DiffSign;				// 등락부호
	char Diff[6];	char _Diff;						// 등락폭
	char DiffRate[5];	char _DiffRate;				// 등락률
	char High[7];	char _High;						// 고가
    char Low[7];	char _Low;						// 저가
	char OfferPrice[7];	char _OfferPrice;			// 매도 호가
    char BidPrice[7];	char _BidPrice;				// 매수 호가
    char Volume[9];	char _Volume;					// 거래량
    char VsPrevVolRate[6];	char _VsPrevVolRate;	// 거래량전일비
    char DiffVolume[8];	char _DiffVolume;			// 변동거래량
    char TrAmount[9];	char _TrAmount;				// 거래대금
    char Open[7];	char _Open;						// 시가
    char WeightAvgPrice[7];	char _WeightAvgPrice;	// 가중평균가
    char Market[1];	char _Market;					// 장구분 ('0': 코스피, '1': 코스닥)
} Tk8OutBlock;

//----------------------------------------------------------------------//
// 코스피/코스닥 ETF NAV (j0/j1)
//----------------------------------------------------------------------//
typedef struct { // 출력
	char skip[3];									// 앞쪽 3바이트는 패킷유형과 압축구분이므로 skip
	char Code[6];	char _Code;						// 종목코드
	char Time[8];	char _Time;						// 시간 (HH:MM:SS)
	char DiffSign[1];	char _DiffSign;				// 등락부호
	char Diff[9];	char _Diff;						// 등락폭
	char NAV_Current[9];	char _NAV_Current;		// NAV 현재가
	char NAV_Open[9];	char _NAV_Open;				// NAV 시가
	char NAV_High[9];	char _NAV_High;				// NAV 고가
	char NAV_Low[9];	char _NAV_Low;				// NAV 저가
	char TrackErrSign[1];	char _TrackingSign;		// 추적 부호
	char TrackingError[9];	char _TrackingError;	// 추적 오차
	char DivergeSign[1];	char _DivergeSign;		// 괴리율 부호
	char DivergeRate[9];	char _DivergeRate;		// 괴리율
} Tj0j1OutBlock;

//----------------------------------------------------------------------//
// 코스피/코스닥 업종 지수 (u1/k1)
//----------------------------------------------------------------------//
typedef struct { // 출력
    char skip[3];									// 앞쪽 3바이트는 패킷유형과 압축구분이므로 skip
	char SectorCode[2];	char _SectorCode;			// 업종코드
	char Time[8];	char _Time;						// 시간
	char IndexValue[8];	char _IndexValue;			// 지수값
	char DiffSign[1];	char _DiffSign;				// 등락부호
	char Diff[8];	char _Diff;						// 등락폭
	char Volume[8];	char _Volume;					// 거래량
	char TrAmount[8];	char _TrAmount;				// 거래대금
	char Open[8];	char _Open;						// 개장 지수값
	char High[8];	char _High;						// 당일 최고값
	char HighTime[8];	char _HighTime;				// 당일 최고값 시간
	char Low[8];	char _Low;						// 당일 최저값
	char LowTime[8];	char _LowTime;				// 당일 최저값 시간
	char DiffRate[5];	char _DiffRate;				// 지수등락률
	char TrVolRate[5];	char _TrVolRate;			// 거래비중 ???
} Tu1k1OutBlock;

/* 코스피/코스닥 업종코드 참고표
코스피 업종명			코스닥 업종명
00 	KRX 100			01 	코스닥지수
01 	코스피지수			03 	기타서비스
02 	대형주			04 	코스닥 IT
03 	중형주			06 	제조
04 	소형주			07 	건설
05 	음식료품			08 	유통
06 	섬유,의복			10 	운송
07 	종이,목재			11 	금융
08 	화학				12 	통신방송서비스
09 	의약품			13 	IT S/W & SVC
10 	비금속광물			14 IT H/W
11 	철강,금속			15 	음식료,담배
12 	기계				16 	섬유,의류
13 	전기,전자			17 	종이,목재
14 	의료정밀			18 	출판,매체복제
15 	운수장비			19 	화학
16 	유통업			20 	제약
17 	전기가스업			21 	비금속
18 	건설업			22 	금속
19 	운수창고			23 	기계,장비
20 	통신업			24 	일반전기전자
21 	금융업			25 	의료,정밀기기
22 	은행				26 	운송장비,부품
24 	증권				27 	기타 제조
25 	보험				28 	통신서비스
26 	서비스업			29 	방송서비스
27 	제조업			30 	인터넷
28 	코스피 200		31 	디지털컨텐츠
29 	코스피 100		32 	소프트웨어
30 	코스피 50			33 	컴퓨터서비스
32 	코스피 배당		34 	통신장비
39 	KP200 건설기계		35 	정보기기
40 	KP200 조선운송		36 	반도체
41 	KP200 철강소재		37 	IT부품
42 	KP200 에너지화학	38 	KOSDAQ 100
43 	KP200 정보기술		39 	KOSDAQ MID 300
44 	KP200 금융		40 	KOSDAQ SMALL
45 	KP200 생활소비재	43 	코스닥 스타
46 	KP200 경기소비재	44 	오락,문화
47 	동일가중 KP200		45 	프리미어
48 	동일가중 KP100		46 	우량기업부
49 	동일가중 KP50		47 	벤처기업부
　	　				48 	중견기업부
　	　				49 	기술성장기업부  */

//----------------------------------------------------------------------//
// 주식 매도(c8101)
//----------------------------------------------------------------------//
typedef struct {
	char pswd_noz8                        [ 44];	char _pswd_noz8;                          //비밀번호
	char issue_codez6                     [  6];	char _issue_codez6;                       //종목번호
	char order_qtyz12                     [ 12];	char _order_qtyz12;                       //주문수량
	char order_unit_pricez10              [ 10];	char _order_unit_pricez10;                //주문단가
	char trade_typez2                     [  2];	char _trade_typez2;                       // 매매유형
    // 00:보통가, 03:시장가, 05:조건부지정가, 12:최유리지정가, 13:최우선지정가,
    // 31 시간외단일가, 61:장전시간외, 71:장후시간외,
    // C0:IOC보통가 (즉시체결.잔량취소), F0:FOK보통가 (즉시전량체결.전량취소),
    // C3:IOC시장가 (즉시체결.잔량취소), F3:FOK시장가 (즉시전량체결.전량취소),
    // C2:IOC최유리 (즉시체결.잔량취소), F2:FOK최유리 (즉시전량체결.전량취소)
	char shsll_pos_flagz1                 [  1];	char _shsll_pos_flagz1;                   //공매도가능여부
	char trad_pswd_no_1z8                 [ 44];	char _trad_pswd_no_1z8;                   //거래비밀번호1
	char trad_pswd_no_2z8                 [ 44];	char _trad_pswd_no_2z8;                   //거래비밀번호2
} Tc8101InBlock;

typedef struct {
	char order_noz10                      [ 10];	char _order_noz10;                        //주문번호
	char order_qtyz12                     [ 12];	char _order_qtyz12;                       //주문수량
	char order_unit_pricez10              [ 10];	char _order_unit_pricez10;                //주문단가
} Tc8101OutBlock;

//----------------------------------------------------------------------//
// 주식 매수(c8102)
//----------------------------------------------------------------------//
typedef struct {
	char pswd_noz8                        [ 44];	char _pswd_noz8;                          //비밀번호
	char issue_codez6                     [  6];	char _issue_codez6;                       //종목번호
	char order_qtyz12                     [ 12];	char _order_qtyz12;                       //주문수량
	char order_unit_pricez10              [ 10];	char _order_unit_pricez10;                //주문단가
	char trade_typez2                     [  2];	char _trade_typez2;                       // 매매유형
    // 00:보통가, 03:시장가, 05:조건부지정가, 12:최유리지정가, 13:최우선지정가,
    // 31 시간외단일가, 61:장전시간외, 71:장후시간외,
    // C0:IOC보통가 (즉시체결.잔량취소), F0:FOK보통가 (즉시전량체결.전량취소),
    // C3:IOC시장가 (즉시체결.잔량취소), F3:FOK시장가 (즉시전량체결.전량취소),
    // C2:IOC최유리 (즉시체결.잔량취소), F2:FOK최유리 (즉시전량체결.전량취소)
	char trad_pswd_no_1z8                 [ 44];	char _trad_pswd_no_1z8;                   //거래비밀번호1
	char trad_pswd_no_2z8                 [ 44];	char _trad_pswd_no_2z8;                   //거래비밀번호2
} Tc8102InBlock;

typedef struct {
	char order_noz10                      [ 10];	char _order_noz10;                        //주문번호
	char order_qtyz12                     [ 12];	char _order_qtyz12;                       //주문수량
	char order_unit_pricez10              [ 10];	char _order_unit_pricez10;                //주문단가
} Tc8102OutBlock;

//----------------------------------------------------------------------//
// 주식 정정 주문 (c8103)
//----------------------------------------------------------------------//
typedef struct {
	char pswd_noz8                        [ 44];	char _pswd_noz8;                          //비밀번호
	char issue_codez6                     [  6];	char _issue_codez6;                       //종목번호
	char crctn_qtyz12                     [ 12];	char _crctn_qtyz12;                       //정정수량
	char crctn_pricez10                   [ 10];	char _crctn_pricez10;                     //정정단가
	char orgnl_order_noz10                [ 10];	char _orgnl_order_noz10;                  //원주문번호
	char all_part_typez1                  [  1];	char _all_part_typez1;                    //정정구분
	char trad_pswd_no_1z8                 [ 44];	char _trad_pswd_no_1z8;                   //거래비밀1
	char trad_pswd_no_2z8                 [ 44];	char _trad_pswd_no_2z8;                   //거래비밀2
} Tc8103InBlock;

typedef struct {
	char orgnl_order_noz10                [ 10];	char _orgnl_order_noz10;                  //원주문번호
	char order_noz10                      [ 10];	char _order_noz10;                        //주문번호
	char mom_order_noz10                  [ 10];	char _mom_order_noz10;                    //모주문번호
	char issue_codez6                     [  6];	char _issue_codez6;                       //후종목번호
	char crctn_qtyz12                     [ 12];	char _crctn_qtyz12;                       //정정수량
	char crctn_pricez10                   [ 10];	char _crctn_pricez10;                     //정정단가
} Tc8103OutBlock;

//----------------------------------------------------------------------//
// 주식 취소 주문 (c8104)
//----------------------------------------------------------------------//
typedef struct {
	char pswd_noz8                        [ 44];	char _pswd_noz8;                          //비밀번호
	char issue_codez6                     [  6];	char _issue_codez6;                       //종목번호
	char canc_qtyz12                      [ 12];	char _canc_qtyz12;                        //취소수량
	char orgnl_order_noz10                [ 10];	char _orgnl_order_noz10;                  //원주문번호
	char all_part_typez1                  [  1];	char _all_part_typez1;                    //취소구분
	char trad_pswd_no_1z8                 [ 44];	char _trad_pswd_no_1z8;                   //거래비밀번호1
	char trad_pswd_no_2z8                 [ 44];	char _trad_pswd_no_2z8;                   //거래비밀번호2
} Tc8104InBlock;

typedef struct {
	char orgnl_order_noz10                [ 10];	char _orgnl_order_noz10;                  //원주문번호
	char order_noz10                      [ 10];	char _order_noz10;                        //주문번호
	char mom_order_noz10                  [ 10];	char _mom_order_noz10;                    //모주문번호
	char issue_codez6                     [  6];	char _issue_codez6;                       //후종목번호
	char canc_qtyz12                      [ 12];	char _canc_qtyz12;                        //취소수량
} Tc8104OutBlock;

//----------------------------------------------------------------------//
// 주식 잔고 조회 (c8201)
//----------------------------------------------------------------------//
typedef struct {
	char pswd_noz44                       [ 44];	char _pswd_noz44;                         //비밀번호
	char bnc_bse_cdz1                     [  1];	char _bnc_bse_cdz1;                       //잔고구분
	// 1:체결기준(주식관련 총평가 -> 기본 )
    // 2:결제잔고
    // 3:시간외종가 체결잔고
    // 4:시간외종가 결제잔고
    // 5:주식잔고평가(주식만 평가)
} Tc8201InBlock;

typedef struct {
	char dpsit_amtz16                     [ 16];	char _dpsit_amtz16;                       //예수금
	char mrgn_amtz16                      [ 16];	char _mrgn_amtz16;                        //신용융자금
	char mgint_npaid_amtz16               [ 16];	char _mgint_npaid_amtz16;                 //이자미납금
	char chgm_pos_amtz16                  [ 16];	char _chgm_pos_amtz16;                    //출금가능금액
	char cash_mrgn_amtz16                 [ 16];	char _cash_mrgn_amtz16;                   //현금증거금
	char subst_mgamt_amtz16               [ 16];	char _subst_mgamt_amtz16;                 //대용증거금
	char coltr_ratez6                     [  6];	char _coltr_ratez6;                       //담보비율
	char rcble_amtz16                     [ 16];	char _rcble_amtz16;                       //현금미수금
	char order_pos_csamtz16               [ 16];	char _order_pos_csamtz16;                 //주문가능액
	char ecn_pos_csamtz16                 [ 16];	char _ecn_pos_csamtz16;                   //ECN주문가능액
	char nordm_loan_amtz16                [ 16];	char _nordm_loan_amtz16;                  //미상환금
	char etc_lend_amtz16                  [ 16];	char _etc_lend_amtz16;                    //기타대여금
	char subst_amtz16                     [ 16];	char _subst_amtz16;                       //대용금액
	char sln_sale_amtz16                  [ 16];	char _sln_sale_amtz16;                    //대주담보금
	char bal_buy_ttamtz16                 [ 16];	char _bal_buy_ttamtz16;                   // 매입원가(계좌합산) 매입가 * 보유수량
	char bal_ass_ttamtz16                 [ 16];	char _bal_ass_ttamtz16;                   // 평가금액(계좌합산) 현재가 * 보유수량
	char asset_tot_amtz16                 [ 16];	char _asset_tot_amtz16;                   // 자산 합계(계좌합산) 평가금액 + 예수금
	char actvt_type10                     [ 10];	char _actvt_type10;                       // 활동유형 (활동, 통합, 폐쇄, 신용계좌)
	char lend_amtz16                      [ 16];	char _lend_amtz16;                        //대출금
	char accnt_mgamt_ratez6               [  6];	char _accnt_mgamt_ratez6;                 //계좌증거금율
	char sl_mrgn_amtz16                   [ 16];	char _sl_mrgn_amtz16;                     //매도증거금
	char pos_csamt1z16                    [ 16];	char _pos_csamt1z16;                      //20%주문가능금액
	char pos_csamt2z16                    [ 16];	char _pos_csamt2z16;                      //30%주문가능금액
	char pos_csamt3z16                    [ 16];	char _pos_csamt3z16;                      //40%주문가능금액
	char pos_csamt4z16                    [ 16];	char _pos_csamt4z16;                      //100%주문가능금액
	char dpsit_amtz_d1_16                 [ 16];	char _dpsit_amtz_d1_16;                   //D1예수금
	char dpsit_amtz_d2_16                 [ 16];	char _dpsit_amtz_d2_16;                   //D2예수금
	char noticez30                        [ 30];	char _noticez30;                          // 공지사항 (사용되지 않음)
	char tot_eal_plsz18                   [ 18];	char _tot_eal_plsz18;                     //총평가손익
	char pft_rtz15                        [ 15];	char _pft_rtz15;                          //수익율
} Tc8201OutBlock;

typedef struct {
	char issue_codez6                     [  6];	char _issue_codez6;                       //종목번호
	char issue_namez40                    [ 40];	char _issue_namez40;                      //종목명
	char bal_typez6                       [  6];	char _bal_typez6;                         //잔고유형
	char loan_datez10                     [ 10];	char _loan_datez10;                       //대출일
	char bal_qtyz16                       [ 16];	char _bal_qtyz16;                         //잔고수량
	char unstl_qtyz16                     [ 16];	char _unstl_qtyz16;                       //미결제량
	char slby_amtz16                      [ 16];	char _slby_amtz16;                        //평균매입가
	char prsnt_pricez16                   [ 16];	char _prsnt_pricez16;                     //현재가
	char lsnpf_amtz16                     [ 16];	char _lsnpf_amtz16;                       //손익(천원)
	char earn_ratez9                      [  9];	char _earn_ratez9;                        //손익율
	char mrgn_codez4                      [  4];	char _mrgn_codez4;                        // 신용유형 (자융,유융,보통,매입)
	char jan_qtyz16                       [ 16];	char _jan_qtyz16;                         //잔량
	char expr_datez10                     [ 10];	char _expr_datez10;                       //만기일
	char ass_amtz16                       [ 16];	char _ass_amtz16;                         //평가금액
	char issue_mgamt_ratez6               [  6];	char _issue_mgamt_ratez6;                 //종목증거금율         /*float->char*/
	char medo_slby_amtz16                 [ 16];	char _medo_slby_amtz16;                   //평균매도가
	char post_lsnpf_amtz16                [ 16];	char _post_lsnpf_amtz16;                  //매도손익
} Tc8201OutBlock1;

//----------------------------------------------------------------------//
// 주문/체결 조회 (s8120)
//----------------------------------------------------------------------//
typedef struct {
	char inq_gubunz1                      [  1];	char _inq_gubunz1;                        // 조회주체구분 (3.계좌별조회 )
	char pswd_noz8                        [ 44];	char _pswd_noz8;                          //비밀번호
	char group_noz4                       [  4];	char _group_noz4;                         // 그룹번호 (0000)
	char mkt_slctz1                       [  1];	char _mkt_slctz1;                         // 시장구분 (0:전체, 1:3일주문, 2:장내채권, 3:제3시장, 4:선물옵션, 5:장외단주,  7:주식옵션현물)
	char order_datez8                     [  8];	char _order_datez8;                       //주문일자
	char issue_codez12                    [ 12];	char _issue_codez12;                      //종목번호
	char comm_order_typez2                [  2];	char _comm_order_typez2;                  // 매체구분 (CC:전체,  AA:영업, BB:온라인)
	char conc_gubunz1                     [  1];	char _conc_gubunz1;                       // 체결구분 (0:전체, 1:미체결, 2:체결)
	char inq_seq_gubunz1                  [  1];	char _inq_seq_gubunz1;                    // 조회순서 (0:번호, 1:모주문번호)
	char sort_gubunz1                     [  1];	char _sort_gubunz1;                       // 정렬구분 (0:주문번호순, 1:주문번호 역순)
	char sell_buy_typez1                  [  1];	char _sell_buy_typez1;                    // 매수도구분 (1:매도, 2:매수, 3:전매, 4:환매)
	char mrgn_typez1                      [  1];	char _mrgn_typez1;                        // 신용구분 (0:보통, 1:신용, 2:대출)
	char accnt_admin_typez1               [  1];	char _accnt_admin_typez1;                 // 계좌구분 (0:전체, 시장구분(mkt_slctz1 = '4')일때 0:전체, 1:지수선물옵션, 2:주식옵션)
	char order_noz10                      [ 10];	char _order_noz10;                        //주문번호
	char ctsz56                           [ 56];	char _ctsz56;                             // CTS (연속처리를 하기위한 Key값??)
	char trad_pswd1z8                     [ 44];	char _trad_pswd1z8;                       //거래비밀번호1
	char trad_pswd2z8                     [ 44];	char _trad_pswd2z8;                       //거래비밀번호2
	char IsPageUp                         [  1];	char _IsPageUp;                           // ISPAGEUP (다음화면이 있는경우:'N', 없는경우:' ')
} Ts8120InBlock;

typedef struct {
	char emp_kor_namez20                  [ 20];	char _emp_kor_namez20;                    //한글사원성명
	char brch_namez30                     [ 30];	char _brch_namez30;                       //한글지점명
	char buy_conc_qtyz14                  [ 14];	char _buy_conc_qtyz14;                    //매수체결수량
	char buy_conc_amtz19                  [ 19];	char _buy_conc_amtz19;                    //매수체결금액
	char sell_conc_qtyz14                 [ 14];	char _sell_conc_qtyz14;                   //매도체결수량
	char sell_conc_amtz19                 [ 19];	char _sell_conc_amtz19;                   //매도체결금액
} Ts8120OutBlock;

typedef struct {
	char order_datez8                     [  8];	char _order_datez8;                       //주문일자
	char order_noz10                      [ 10];	char _order_noz10;                        //주문번호
	char orgnl_order_noz10                [ 10];	char _orgnl_order_noz10;                  //원주문번호
	char accnt_noz11                      [ 11];	char _accnt_noz11;                        //계좌번호
	char accnt_namez20                    [ 20];	char _accnt_namez20;                      //계좌명
	char order_kindz20                    [ 20];	char _order_kindz20;                      // 주문구분 ('매수|매도' +[' ' + (정정|취소)] ex) 현금매수, 현금매수 정정, 현금매수 취소)
	char trd_gubun_noz1                   [  1];	char _trd_gubun_noz1;                     // 매매구분번호 (더 이상 사용되지 않는 듯)
	char trd_gubunz20                     [ 20];	char _trd_gubunz20;                       // 매매구분 (보통,조건부지정,시장가,장전)
	char trade_type_noz1                  [  1];	char _trade_type_noz1;                    // 거래구분번호 (더 이상 사용되지 않는 듯)
	char trade_type1z20                   [ 20];	char _trade_type1z20;                     // 거래구분  (더 이상 사용되지 않는 듯)
	char issue_codez12                    [ 12];	char _issue_codez12;                      //종목번호
	char issue_namez40                    [ 40];	char _issue_namez40;                      //종목명
	char order_qtyz10                     [ 10];	char _order_qtyz10;                       //주문수량
	char conc_qtyz10                      [ 10];	char _conc_qtyz10;                        //체결수량
	char order_unit_pricez12              [ 12];	char _order_unit_pricez12;                //주문단가
	char conc_unit_pricez12               [ 12];	char _conc_unit_pricez12;                 //체결평균단가
	char crctn_canc_qtyz10                [ 10];	char _crctn_canc_qtyz10;                  // 정정취소수량 (정정/취소 된? 정정/취소 후?)
	char cfirm_qtyz10                     [ 10];	char _cfirm_qtyz10;                       //확인수량
	char media_namez12                    [ 12];	char _media_namez12;                      //매체구분
	char proc_emp_noz5                    [  5];	char _proc_emp_noz5;                      //처리사번
	char proc_timez8                      [  8];	char _proc_timez8;                        //처리시간
	char proc_termz8                      [  8];	char _proc_termz8;                        //처리단말
	char proc_typez12                     [ 12];	char _proc_typez12;                       // 처리구분 ('정상', 정정/취소인경우 '확인')
	char rejec_codez5                     [  5];	char _rejec_codez5;                       // 거부코드 (거부시만 코드제공)
	char avail_qtyz10                     [ 10];	char _avail_qtyz10;                       //정취가능수량
	char mkt_typez1                       [  1];	char _mkt_typez1;                         //시장구분
	char shsll_typez20                    [ 20];	char _shsll_typez20;                      // 공매도구분 ('정상')
	char passwd_noz8                      [  8];	char _passwd_noz8;                        // 비밀번호 (??비밀번호가 왜 8자리??)
} Ts8120OutBlock1;

typedef struct {
	char ctsz56                           [ 56];	char _ctsz56;                             //CTS
	char nextbutton                       [  1];	char _nextbutton;                         //NEXTBUTTON
} Ts8120OutBlock_IN;

//----------------------------------------------------------------------//
// 매도 가능 수량 (p8101)
//----------------------------------------------------------------------//
typedef struct tagp8101InBlock {
	char pswd_noz8                        [ 44];	char _pswd_noz8;                          //비밀번호
	char gubunz1                          [  1];	char _gubunz1;                            //구분
} Tp8101InBlock;

typedef struct tagp8101OutBlock {
	char accnt_namez30                    [ 30];	char _accnt_namez30;                      //계좌명               /*신OBM에존재하지않는항목*/
} Tp8101OutBlock;

typedef struct tagp8101OutBlock1 {
	char gubunz1                          [  1];	char _gubunz1;                            // 구분
                                            // 1:현금, 2:융자, 3:채권, 4:대주, 5:대출주식, 6:융자주식합계, 7:대출주식합계,
                                            // 8:융자주식 및 대출주식, 9:융자주식합계+대출주식합계, A:전체
	char gubun_namez6                     [  6];	char _gubun_namez6;                       // 구분명 (현금, 융자, 대출)
	char issue_codez12                    [ 12];	char _issue_codez12;                      //종목코드
	char issue_namez30                    [ 30];	char _issue_namez30;                      //종목명
	char mrgn_typez10                     [ 10];	char _mrgn_typez10;                       // 신용구분 (유통융자, 자기융자, 매입자금)
	char lend_datez10                     [ 10];	char _lend_datez10;                       //대출일자
	char taxtn_typez10                    [ 10];	char _taxtn_typez10;                      // 과세유형 (더 이상 사용되지 않는 듯)
	char bal_qtyz12                       [ 12];	char _bal_qtyz12;                         //잔고수량
	char sell_rcble_qtyz12                [ 12];	char _sell_rcble_qtyz12;                  //매도미결제
	char buy_rcble_qtyz12                 [ 12];	char _buy_rcble_qtyz12;                   //매수미결제
	char sell_psqtyz12                    [ 12];	char _sell_psqtyz12;                      //매도가능수량
	char today_sell_rcble_qz12            [ 12];	char _today_sell_rcble_qz12;              //당일매도미체결수량
	char avrg_purch_uprc                  [ 10];	char _avrg_purch_uprc;                    //매입단가
} Tp8101OutBlock1;

//----------------------------------------------------------------------//
// 매수 가능 수량 (p8105)
//----------------------------------------------------------------------//
typedef struct tagp8105InBlock {
	char pwdz8                            [ 44];	char _pwdz8;                              //비밀번호
	char ost_dit_cdz1                     [  1];	char _ost_dit_cdz1;                       //구분코드             /*1현금2:신용3:매입자금대출*/
	char sby_dit_cdz1                     [  1];	char _sby_dit_cdz1;                       //매매구분코드         /*1:매도상환2:매수신규*/
	char iem_gbz1                         [  1];	char _iem_gbz1;                           //종목구분             /*1:주식2:ELW3:신주인수4:기타*/
	char iem_cdz12                        [ 12];	char _iem_cdz12;                          //종목코드
	char nmn_pr_tp_gbz1                   [  1];	char _nmn_pr_tp_gbz1;                     //호가유형구분         /*1:구호가구분-1자리2:신시스템호가구분-2자리*/
	char nmn_pr_tp_cdz2                   [  2];	char _nmn_pr_tp_cdz2;                     //호가유형코드
	// 01:보통, 05:시장가, 06:조건부 지정가, 10:S-OPTION자기, 11:금전신탁, 12:최유리지정가, 13:최우선지정가,
    // 61:장전 시간외, 71:장후 시간외 종가, 81:시간외 단일가
	char orr_prz18                        [ 18];	char _orr_prz18;                          //주문가격
	char mdi_tp_cdz1                      [  1];	char _mdi_tp_cdz1;                        //매체유형코드
	// 1:지점, 2:HTS/인터넷, 3:모바일/ARS, 4:고객지원센터, 5:TX flat, 6:TX Leverage(일반),
	// 7:TX Lever(대출/신용), 8:TX Winner(일반), 9:TX Winner(STOP-MIT), 0:TX 바로주문
	char cfd_lon_cdz2                     [  2];	char _cfd_lon_cdz2;                       // 신용대출코드 (01:유통융자, 02:자기융자, 03:유통대주, 04:자기대주)
	char lon_dtz8                         [  8];	char _lon_dtz8;                           //대출일자
} Tp8105InBlock;

typedef struct {
	char dcaz18                           [ 18];	char _dcaz18;                             //예수금               /*금일예수금*/
	char nxt_dd_dcaz18                    [ 18];	char _nxt_dd_dcaz18;                      //익일예수금           /*D+1예수금*/
	char nxt2_dd_dcaz18                   [ 18];	char _nxt2_dd_dcaz18;                     //익익일예수금         /*D+2예수금*/
	char max_pbl_amtz18                   [ 18];	char _max_pbl_amtz18;                     //최대가능금액         /*미수가능금액*/
	char max_pbl_qtyz18                   [ 18];	char _max_pbl_qtyz18;                     //최대가능수량         /*미수가능수량*/
	char rvb_orn_max_pbl_feez18           [ 18];	char _rvb_orn_max_pbl_feez18;             //미수발생최대가능수수료 /*미수수수료*/
	char csh_orr_pbl_amtz18               [ 18];	char _csh_orr_pbl_amtz18;                 //현금주문가능금액     /*현금가능금액*/
	char csh_orr_pbl_qtyz18               [ 18];	char _csh_orr_pbl_qtyz18;                 //현금주문가능수량     /*현금가능수량*/
	char ost_fee1z18                      [ 18];	char _ost_fee1z18;                        //현금수수료           /*현금수수료*/
	char cfd_rvb_orr_pbl_amtz18           [ 18];	char _cfd_rvb_orr_pbl_amtz18;             //신용미수주문가능금액 /*신용미수가능금액*/
	char cfd_rvb_orr_pbl_qtyz18           [ 18];	char _cfd_rvb_orr_pbl_qtyz18;             //신용미수주문가능수량 /*신용미수가능수량*/
	char cfd_max_pbl_feez18               [ 18];	char _cfd_max_pbl_feez18;                 //신용최대가능수수료   /*신용미수수수료*/
	char cfd_orr_pbl_amtz18               [ 18];	char _cfd_orr_pbl_amtz18;                 //신용주문가능금액     /*신용미발생가능금액*/
	char cfd_orr_pbl_qtyz18               [ 18];	char _cfd_orr_pbl_qtyz18;                 //신용주문가능수량     /*신용미발생가능수량*/
	char ost_fee2z18                      [ 18];	char _ost_fee2z18;                        //수수료2              /*신용미발생수수료*/
	char sdr_xps1z18                      [ 18];	char _sdr_xps1z18;                        //제비용1
	char sdr_xpsz18                       [ 18];	char _sdr_xpsz18;                         //제비용
} Tp8105OutBlock;

//----------------------------------------------------------------------//
// 개별 주식 매도 수량 (p8104) : 현금 주식 (??)
//----------------------------------------------------------------------//
typedef struct {
	char pswd_noz8                        [ 44];	char _pswd_noz8;                          //비밀번호
	char issue_codez6                     [  6];	char _issue_codez6;                       //종목코드
	char gubunz1                          [  1];	char _gubunz1;                            // 구분 (1:현금, 2:신용(유통융자))
	char the_datez8                       [  8];	char _the_datez8;                         //대출일               /*신OBM에존재하지않는항목*/
} Tp8104InBlock;

typedef struct {
	char issue_codez6                     [  6];	char _issue_codez6;                       //종목코드
	char order_qtyz12                     [ 12];	char _order_qtyz12;                       //매도가능수량
} Tp8104OutBlock;

//----------------------------------------------------------------------//
// 체결 확인 실시간 패킷 (d2) - 거래소 발송
//----------------------------------------------------------------------//
typedef struct {
    char skip[3];				// 앞쪽 3바이트는 패킷유형과 압축구분이므로 skip
	char userid                           [  8];   //사용자ID
	char itemgb                           [  1];   //ITEM구분 (1:주식, 2:장내선물, 4:주식선물)
	char accountno                        [ 11];   //계좌번호
	char orderno                          [ 10];   //주문번호
	char issuecd                          [ 15];   //종목코드
	char slbygb                           [  1];   // 매도수구분 (1:매도, 2:매수, 3:전매, 4:환매)
	char concgty                          [ 10];   //체결수량
	char concprc                          [ 11];   //체결가격
	char conctime                         [  6];   //체결시간
	char ucgb                             [  1];   // 정정취소구분 (0:체결, 1:정정, 2:취소, 3:주문거부, 4:IOC취소, 5:FOK취소)
	char rejgb                            [  1];   // 거부구분 (0:정상, 1:거부)
	char fundcode                         [  3];   // 펀드코드
	char sin_gb                           [  2];   // 신용구분
    // 10:정상, 21:자기융자매수, 22:자기융자매도상환, 23:자기대주매도, 24:자기대주매수상환,
    // 31:유통융자매수, 32:유통융자매도상환, 33:유통대주매도, 34:유통대주매수상환,
    // 61:청약대출매도, 62:보통대출매도, 63:매입대출매수, 64:매입대출매도
	char loan_date                        [  8];   //대출일자
	char ato_ord_tpe_chg                  [  1];   //선물옵션주문유형변경여부
	char filler                           [ 34];   //filler
} Td2OutBlock;

//----------------------------------------------------------------------//
// 주문 확인 실시간 패킷 (d3) - 증권사 발송
//----------------------------------------------------------------------//
typedef struct {
    char skip[3];				// 앞쪽 3바이트는 패킷유형과 압축구분이므로 skip
	char userid                           [  8];   //USERID
	char itemgb                           [  1];   // ITEM구분 (1:주식, 2:장내선물)
	char accountno                        [ 11];   //계좌번호
	char orderno                          [ 10];   //주문번호
	char orgordno                         [ 10];   //원주문번호
	char ordercd                          [  2];   // 주문구분 (12:정정, 13:취소)
	char issuecd                          [ 15];   //종목코드
	char issuename                        [ 20];   //종목명
	char slbygb                           [  1];   // 매매구분 (1:매도, 2:매수)
	char order_type                       [  2];   //주문유형
	// 주식 : 01~81,
    // 01:지정가, 05:시장가, 06:조건부지정가, 09:자기주식, 10:S-OPTION자사주,
    // 11:금전신탁자기주식, 12:최유리지정가, 13 최우선지정가,
    // 51:장중대량, 52:장중바스켓,
    // 61:장전시간외종가, 62:장전시간외대량, 63:장전바스켓,
    // 67:장전대량신탁자사주, 69:장전대량자기주식
    // 71:장후시간외종가, 72:장후시간외대량, 77:신탁자기주식시간외대량, 79:자기주식시간외대량,
    // 80:바스켓매매, 81 시간외단일가,
    // 선물/옵션/주식선물/주식옵션 : B,C,L,M
    // B:최유리, C:조건부, L:지정가, M:시장가
	char ordergty                         [ 10];   //주문수량
	char orderprc                         [ 11];   //주문단가
	char procnm                           [  2];   //처리구분
	char commcd                           [  2];   //매체구분
	char order_cond                       [  1];   // 주문조건 (1:IOC, 2:FOK)
	char fundcode                         [  3];   //펀드코드
	char sin_gb                           [  2];   // 신용구분
    // 10:정상, 21:자기융자매수, 22:자기융자매도상환, 23:자기대주매도, 24:자기대주매수상환,
    // 31:유통융자매수, 32:유통융자매도상환, 33:유통대주매도, 34:유통대주매수상환,
    // 61:청약대출매도, 62:보통대출매도, 63:매입대출매수, 64:매입대출매도
	char order_time                       [  6];   //주문시간
	char loan_date                        [  8];   //대출일자
} Td3OutBlock;