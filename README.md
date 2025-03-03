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

보통 일일이 `t.Run`을 작성하기보다, 테스트 케이스에 대한 struct를 미리 만들어놓고 반복문을 돌리는 관례가 있는 듯 하다.(레퍼런스 확인 필요)

또한 찾아봤는데 기본적으로 `t.Group`과 같은 그룹핑이 없다. 테스트 함수 자체를 하나의 그룹으로 보고 하위에 `t.Run`을 여러개 두는 듯?

### 슬라이스가 동일한지 확인하기 위해서는 `==`가 아닌 `reflect.DeepEqual`을 사용.

이건 JS에서도 비슷한데, slice 자체는 내부의 값이 아닌 메모리 주소로 비교를 하는 듯? JS에서는 `JSON.stringify`로 문자열로 변환하여 확인하거나, 별도 패키지(`lodash` 등)을 사용해서 비교하는데 Go에서는 내부 패키지로 `reflect.DeepEqual`을 제공한다.
