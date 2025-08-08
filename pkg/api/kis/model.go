package kis

// KISBalanceResponse 한국투자증권 해외주식 잔고 조회 응답
type KISBalanceResponse struct {
	RtCd         string              `json:"rt_cd"`
	MsgCd        string              `json:"msg_cd"`
	Msg1         string              `json:"msg1"`
	CtxAreaFk200 string              `json:"ctx_area_fk200"`
	CtxAreaNk200 string              `json:"ctx_area_nk200"`
	Output1      []KISBalanceOutput1 `json:"output1"`
	Output2      KISBalanceOutput2   `json:"output2"`
}

// KISBalanceOutput1 개별 종목 잔고 정보 (해외주식)
type KISBalanceOutput1 struct {
	Cano            string `json:"cano"`               // 종합계좌번호
	AcntPrdtCd      string `json:"acnt_prdt_cd"`       // 계좌상품코드
	PrdtTypeCd      string `json:"prdt_type_cd"`       // 상품유형코드
	OvrsPdno        string `json:"ovrs_pdno"`          // 해외상품번호
	OvrsItemName    string `json:"ovrs_item_name"`     // 해외종목명
	FrcrEvluPflsAmt string `json:"frcr_evlu_pfls_amt"` // 외화평가손익금액
	EvluPflsRt      string `json:"evlu_pfls_rt"`       // 평가손익율
	PchsAvgPric     string `json:"pchs_avg_pric"`      // 매입평균가격
	OvrsCblcQty     string `json:"ovrs_cblc_qty"`      // 해외잔고수량
	OrdPsblQty      string `json:"ord_psbl_qty"`       // 주문가능수량
	FrcrPchsAmt1    string `json:"frcr_pchs_amt1"`     // 외화매입금액1
	OvrsStckEvluAmt string `json:"ovrs_stck_evlu_amt"` // 해외주식평가금액
	NowPric2        string `json:"now_pric2"`          // 현재가격2
	TrCrcyCd        string `json:"tr_crcy_cd"`         // 거래통화코드
	OvrsExcgCd      string `json:"ovrs_excg_cd"`       // 해외거래소코드
	LoanTypeCd      string `json:"loan_type_cd"`       // 대출유형코드
	LoanDt          string `json:"loan_dt"`            // 대출일자
	ExpdDt          string `json:"expd_dt"`            // 만기일자
}

// KISBalanceOutput2 전체 포트폴리오 요약 정보
type KISBalanceOutput2 struct {
	FrcrPchsAmt1     string `json:"frcr_pchs_amt1"`      // 외화매입금액1
	OvrsRlztPflsAmt  string `json:"ovrs_rlzt_pfls_amt"`  // 해외실현손익금액
	OvrsTotPfls      string `json:"ovrs_tot_pfls"`       // 해외총손익
	RlztErngRt       string `json:"rlzt_erng_rt"`        // 실현수익율
	TotEvluPflsAmt   string `json:"tot_evlu_pfls_amt"`   // 총평가손익금액
	TotPftrt         string `json:"tot_pftrt"`           // 총수익률
	FrcrBuyAmtSmtl1  string `json:"frcr_buy_amt_smtl1"`  // 외화매수금액합계1
	OvrsRlztPflsAmt2 string `json:"ovrs_rlzt_pfls_amt2"` // 해외실현손익금액2
	FrcrBuyAmtSmtl2  string `json:"frcr_buy_amt_smtl2"`  // 외화매수금액합계2
}

// KISPriceResponse 한국투자증권 해외주식 현재가 조회 응답
type KISPriceResponse struct {
	RtCd   string         `json:"rt_cd"`
	MsgCd  string         `json:"msg_cd"`
	Msg1   string         `json:"msg1"`
	Output KISPriceOutput `json:"output"`
}

// KISPriceOutput 해외주식 현재가 정보
type KISPriceOutput struct {
	Rsym   string `json:"rsym"`    // 실시간조회종목코드
	Pvol   string `json:"pvol"`    // 전일거래량
	Open   string `json:"open"`    // 시가
	High   string `json:"high"`    // 고가
	Low    string `json:"low"`     // 저가
	Last   string `json:"last"`    // 현재가
	Base   string `json:"base"`    // 전일종가
	Tomv   string `json:"tomv"`    // 시가총액
	Pamt   string `json:"pamt"`    // 전일거래대금
	Uplp   string `json:"uplp"`    // 상한가
	Dnlp   string `json:"dnlp"`    // 하한가
	H52p   string `json:"h52p"`    // 52주최고가
	H52d   string `json:"h52d"`    // 52주최고일자
	L52p   string `json:"l52p"`    // 52주최저가
	L52d   string `json:"l52d"`    // 52주최저일자
	Perx   string `json:"perx"`    // PER
	Pbrx   string `json:"pbrx"`    // PBR
	Epsx   string `json:"epsx"`    // EPS
	Bpsx   string `json:"bpsx"`    // BPS
	Shar   string `json:"shar"`    // 상장주수
	Mcap   string `json:"mcap"`    // 자본금
	Curr   string `json:"curr"`    // 통화
	Zdiv   string `json:"zdiv"`    // 소수점자리수
	Vnit   string `json:"vnit"`    // 매매단위
	TXprc  string `json:"t_xprc"`  // 원환산당일가격
	TXdif  string `json:"t_xdif"`  // 원환산당일대비
	TXrat  string `json:"t_xrat"`  // 원환산당일등락
	PXprc  string `json:"p_xprc"`  // 원환산전일가격
	PXdif  string `json:"p_xdif"`  // 원환산전일대비
	PXrat  string `json:"p_xrat"`  // 원환산전일등락
	TRate  string `json:"t_rate"`  // 당일환율
	PRate  string `json:"p_rate"`  // 전일환율
	TXsgn  string `json:"t_xsgn"`  // 원환산당일기호
	PXsng  string `json:"p_xsng"`  // 원환산전일기호
	EOrdyn string `json:"e_ordyn"` // 거래가능여부
	EHogau string `json:"e_hogau"` // 호가단위
	EIcod  string `json:"e_icod"`  // 업종(섹터)
	EParp  string `json:"e_parp"`  // 액면가
	Tvol   string `json:"tvol"`    // 거래량
	Tamt   string `json:"tamt"`    // 거래대금
	EtypNm string `json:"etyp_nm"` // ETP 분류명
}
