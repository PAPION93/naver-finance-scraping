# Naver Finance Scraping

> 네이버 금융에서 종목별 데이터를 스크래핑한다.

## Scrapping list

### 기관 매매동향

> 기관 매수 동향 파악하여 리스트화  
> 가장 최근 날짜의 기관의 매수량 8000 이상 상승하였을 경우만  
> 종가 10,000 이상  
> 전일비 하락

## About Go

### 1. goquery

    go get github.com/PuerkitoBio/goquery

### 2. Echo Framework

## To-do

1. 다음 금융으로 수정
1. 해당종목 당일 매매현황만 가져올 것.

1. 외국인 지분율에 따라 외국인을 참고할 것.
    > 지분율 기준치 정할 것  
    > 전날대비 상승률로 확인해볼 것

1. 시각성 높이기
1. 카카오톡 챗봇으로 받고싶은 데이터
    - 저장한 파일 기반한 추천 종목
    - 종목 검색
