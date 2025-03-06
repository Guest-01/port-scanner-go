# port-scanner-go

> 간단한 TCP 포트 스캐너

Go를 학습하기 위한 프로젝트로, 외부 패키지를 쓰지 않고 간단한 CLI를 만들어보기.

## 알게된 사실들

### 문자열 포맷팅은 함수를 사용해야한다

Go에서는 js나 python과 다르게 문자열 자체에 포매팅 기능이 없다. 만일 변수 등을 포함한 문자열을 만들고 싶다면 `fmt` 패키지의 함수들을 사용해야한다. 특히 그대로 출력하고자 한다면 바로 `fmt.Printf` 등을 사용하면 되겠지만, 그게 아니라 먼저 만든 문자열을 변수에 담아두고 싶다면 마찬가지로 `fmt.Sprintf` 함수를 써서 리턴 값으로 문자열을 만들어야한다.

자바스크립트에서는...
```ts
const formatstr = `myVar: ${myVar}`;
```

하지만 고에서는...
```golang
formatstr := fmt.Sprintf("myVar: %v", myVar)
```

### if문에 조건절 앞에 문장이 올 수 있다, 그리고 삼항 연산자가 없다(!)

그리고 앞에 문장에서 선언한 변수는 if문(else 포함) 전체에서 사용 가능하다. 그리고 삼항 연산자가 없어서 아래 예시를 보면 다른건 다 같고, `isOpen` 여부에 따라 한단어만 다르게 하고 싶은건데 `if`-`else`를 써야만 했다...

```golang
if isOpen := scanPort(host, port); isOpen {
    fmt.Printf("%-5d : Open\n", port)
} else {
    fmt.Printf("%-5d : Closed\n", port)
}
```

### 테스트는 함수별로, 그리고 `t.Run`을 통해 여러 케이스를 실행할 수 있다.

보통 일일이 `t.Run`을 작성하기보다, 테스트 케이스에 대한 struct를 미리 만들어놓고 반복문을 돌리는 관례가 있는 듯 하다.~~(레퍼런스 확인 필요)~~ 레퍼런스: [GoWiki - Table Driven Tests](https://go.dev/wiki/TableDrivenTests)

또한 찾아봤는데 기본적으로 `t.Group`과 같은 그룹핑이 없다. 테스트 함수 자체를 하나의 그룹으로 보고 하위에 `t.Run`을 여러개 두는 듯?

### 슬라이스가 동일한지 확인하기 위해서는 `==`가 아닌 `reflect.DeepEqual`을 사용.

이건 JS에서도 비슷한데, slice 자체는 내부의 값이 아닌 메모리 주소로 비교를 하는 듯? JS에서는 `JSON.stringify`로 문자열로 변환하여 확인하거나, 별도 패키지(`lodash` 등)을 사용해서 비교하는데 Go에서는 내부 패키지로 `reflect.DeepEqual`을 제공한다.

### Go에서 에러를 만드는 법은 세가지가 있다.

1) `errors.New`: 정적인 에러 메시지를 출력할 때.
2) `fmt.Errorf`: 동적인 에러 메시지(포매팅) 가능.
3) 커스텀 에러 인터페이스 만들기: 에러를 이름으로 분류를 하고 싶을때(?)

지금은 2번으로 손쉽게 해결했지만 3번에 대해 좀 더 알아볼 필요가 있음.

### Slice와 Slice를 붙이고 싶은 경우 `...`를 사용할 수 있다

Slice에 값을 추가하려는 경우에는 `append()` 함수를 쓴다. 그런데 슬라이스에 슬라이스를 추가하려고하면 어떻게할까 고민했는데, Go에도 Javascript와 비슷하게 Unpacking 같은 기능이 있다. `append()`함수는 두번째 인자부터 개수 제한 없이 받을 수 있는데, 이를 이용해서 슬라이스를 풀어서 넣을 수 있다.

```go
append(slice1, slice2...)
```

JS와 다른 점은 ...이 앞이 아닌 뒤에 붙는다는 것이다.

### Goroutine에서 안전하게 슬라이스 삽입하기

결과를 포트 순으로 출력하기 위해 고루틴을 돌며 바로 출력하지 않고 결과를 저장할 슬라이스를 만들어서 넣은 후에 정렬을 하고 출력하게 되었다. 그런데 보통 이렇게 슬라이스를 고루틴에서 쓰게 되면 race condition이 발생할 수 있다고 하여 권장하지 않고, 대신 `channel`을 쓰거나 일일이 `mutex.Lock()`을 걸어줘야하는 것으로 알고 있었다.

그런데 아래와 같이 정해진 크기만큼의 슬라이스를 미리 만들고, 인덱스를 가지고 접근하면 괜찮다고 한다. (Claude3.7 답변 기준)

```go
results := make([]ScanResult, len(ports))

for i, port := range ports {
    wg.Add(1)
    go func(i, port int) {
        defer wg.Done()
        isOpen := scanPort(host, port)
        results[i] = ScanResult{Port: port, IsOpen: isOpen}
    }(i, port)
}
```

보통 문제는 같은 메모리에 접근하거나 슬라이스를 `append`를 이용하여 넣었을 때 확장이 일어나면서 생긴다고 한다. 위 같은 코드에서는 그런 문제가 일어나지 않기 때문에 괜찮다고 한다. 실제로 테스트를 여러번 해봤는데 별다른 오류가 생기지 않더라.

### TODO

- ~~복합 포트 입력 받기 (예: 80, 443, 8080-8081)~~ 완료.
- ~~결과를 포트순으로 출력하기~~ 완료.