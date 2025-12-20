# Quantum Permutation Pad (量子置换一次一密)

[![GoDoc][1]][2] [![Go Report Card][3]][4] [![CreatedAt][5]][6] 

[1]: https://godoc.org/github.com/xtaci/qpp?status.svg
[2]: https://pkg.go.dev/github.com/xtaci/qpp
[3]: https://goreportcard.com/badge/github.com/xtaci/qpp
[4]: https://goreportcard.com/report/github.com/xtaci/qpp
[5]: https://img.shields.io/github/created-at/xtaci/qpp
[6]: https://img.shields.io/github/created-at/xtaci/qpp

[Quantum Permutation Pad](https://link.springer.com/content/pdf/10.1140/epjqt/s40507-023-00164-3.pdf) (QPP) 是一种旨在利用量子力学原理进行安全通信的加密协议。虽然具体的实现细节可能因理论模型而异，但其核心概念涉及利用量子特性（如叠加和纠缠）来增强数据传输的安全性。本文档概述了 QPP 及其与量子力学和密码学的关系。

## Quantum Permutation Pad 的核心概念

1. **量子力学原理**：QPP 依赖于基本的量子力学原理，特别是叠加（量子比特同时处于多种状态的能力）和纠缠（无论距离多远，量子比特之间都存在关联）。

2. **量子比特 (Qubits)**：与经典比特（只能是 0 或 1）不同，QPP 使用量子比特，它可以处于 0、1 或这两个状态的任意量子叠加态。

3. **置换操作**：在 QPP 的语境中，置换是指**转换/解释**明文的方式，这与经典密码学中重复重新排列密钥不同。置换的总数可以表示为 $P_n$。对于一个 8 位的字节，总的置换数为 $P_{256} =$
```256! =857817775342842654119082271681232625157781520279485619859655650377269452553147589377440291360451408450375885342336584306157196834693696475322289288497426025679637332563368786442675207626794560187968867971521143307702077526646451464709187326100832876325702818980773671781454170250523018608495319068138257481070252817559459476987034665712738139286205234756808218860701203611083152093501947437109101726968262861606263662435022840944191408424615936000000000000000000000000000000000000000000000000000000000000000```
4. **密钥空间扩展**：8 位系统的经典密钥空间有 256 种可能的密钥（整数）。而在 QPP 中，量子密钥空间扩展到了 256! 种可能的置换算子或门！
5. **实现**：QPP 既可以使用矩阵进行经典实现，也可以使用量子门进行量子力学实现。


## 应用与优势

- **高安全性**：通过利用量子力学的独特属性，QPP 提供了优于经典加密方法的安全级别。
- **面向未来**：随着量子计算机变得越来越强大，经典加密方案（如 RSA 和 ECC）面临着越来越多的漏洞。QPP 提供了一种抗量子的替代方案。
- **安全通信**：QPP 非常适合量子网络中的安全通信以及保护高度敏感的数据。

## 使用示例
内部 PRNG（不推荐）
```golang
func main() {
    seed := make([]byte, 32)
    io.ReadFull(rand.Reader, seed)

    qpp := NewQPP(seed, 977) // 这里的 pad 数量是一个质数

    msg := make([]byte, 65536)
    io.ReadFull(rand.Reader, msg)

    qpp.Encrypt(msg)
    qpp.Decrypt(msg)
}
```

外部 PRNG 与共享 pads（**推荐**）
```golang
func main() {
    seed := make([]byte, 32)
    io.ReadFull(rand.Reader, seed)

    qpp := NewQPP(seed, 977)

    msg := make([]byte, 65536)
    io.ReadFull(rand.Reader, msg)

    rand_enc := qpp.CreatePRNG(seed)
    rand_dec := qpp.CreatePRNG(seed)

    qpp.EncryptWithPRNG(msg, rand_enc)
    qpp.DecryptWithPRNG(msg, rand_dec)
}
```

NewQPP 生成的置换如下所示（采用[轮换表示法](https://zh.wikipedia.org/wiki/%E7%BD%AE%E6%8D%A2#%E8%BD%AE%E6%8D%A2%E8%A1%A8%E7%A4%BA%E6%B3%95)）：
```
(0 4 60 108 242 196)(1 168 138 16 197 29 57 21 22 169 37 74 205 33 56 5 10 124 12 40 8 70 18 6 185 137 224)(2 64 216 178 88)(3 14 98 142 128 30 102 44 158 34 72 38 50 68 28 154 46 156 254 41 218 204 161 194 65)(7 157 101 181 141 121 77 228 105 206 193 155 240 47 54 78 110 90 174 52 207 233 248 167 245 199 79 144 162 149 97
140 111 126 170 139 175 119 189 171 215 55 89 81 23 134 106 251 83 15 173 250 147 217 115 229 99 107 223 39 244 246
 225 252 226 203 235 236 253 43 188 209 145 184 91 31 49 84 210 117 59 133 129 75 150 127 200 130 132 247 159 241
255 71 120 63 249 201 212 131 95 222 238 125 237 109 186 213 151 176 143 202 179 232 103 148 191 239)(9 20 113 73 69 160 114 122 164 17 208 58 116 36 26 96 24)(11 80 32 152 146 82 53 62 66 76 86 51 112 221 27 163 180 214 123 219 234)
(13 42 166 25 165)(19 172 177 230 198 45 61 104 136 100 182 85 153 35 192 48 220 94 190 118 195)(67)(87 92 93 227 211)(135)(183 243)(187)(231)
```
![circular](https://github.com/user-attachments/assets/3fa50405-1b4e-4679-a495-548850c4315b)

## 本实现的安全性设计
整体安全性相当于 **1683 位** 对称加密。

8 量子比特系统中的置换矩阵数量由提供的种子决定并随机选择。
<img width="867" alt="image" src="https://github.com/user-attachments/assets/93ce8634-5300-47b1-ba1b-e46d9b46b432">

置换 pad 可以用[轮换表示法](https://zh.wikipedia.org/wiki/%E7%BD%AE%E6%8D%A2#%E8%BD%AE%E6%8D%A2%E8%A1%A8%E7%A4%BA%E6%B3%95)写为： $\sigma =(1\ 2\ 255)(3\ 36)(4\ 82\ 125)(....)$ ，其中的元素不像其他[流密码](https://zh.wikipedia.org/wiki/%E6%B5%81%E5%AF%86%E7%A0%81)那样可以通过两次 **异或 (XOR)** 来还原。

#### 局部性与随机性
随机置换会破坏[局部性原理](https://zh.wikipedia.org/wiki/%E5%B1%80%E9%83%A8%E6%80%A7%E5%8E%9F%E7%90%86)，而这对性能至关重要。为了获得更高的加密速度，必须保持一定程度的局部性。在本设计中，我们不是每个字节都切换 pad，而是每 8 个字节使用一个新的随机 pad。

![349804164-3f6da444-a9f4-4d0a-b190-d59f2dca9f00](https://github.com/user-attachments/assets/2358766e-a0d3-4c21-93cb-c221aa0cece0)

上图清楚地表明，每个字节都切换 pad 会导致性能低下，而每 8 个字节切换一次 pad 则能产生足够的性能。

您可以直接从 https://github.com/xtaci/kcptun/releases 下载并启用 ```-QPP``` 选项来尝试。

#### 性能
在现代 CPU 中，最新的 QPP 优化可以轻松达到超过 1GB/s 的速度。
![348621244-4061d4a9-e7fa-43f5-89ef-f6ef6c00a2e7](https://github.com/user-attachments/assets/78952157-df39-4088-b423-01f45548b782)

## 设置 PADs 的安全注意事项

Pads 的数量理想情况下应与 8 互素 (coprime)，因为结果表明 PRNG 中存在与数字 8 相关的隐藏结构。

![88d8de919445147f5d44ee059cca371](https://github.com/user-attachments/assets/9e1a160d-5433-4e24-9782-2ae88d87453d)

我们使用真实数据（加密后的圣经）演示了使用 64 个 pads 和 15 个 pads 加密圣经的情况。

**对于 Pads(64)**，由于 $GCD(64,8) == 8, \chi^2 =3818 $

![348794146-4f6d5904-2663-46d7-870d-9fd7435df4d0](https://github.com/user-attachments/assets/e2a67bad-7d10-46e4-8d23-9866918ef04b)

**对于 Pads(15)**，由于 $GCD(15,8) == 1,\chi^2 =230$，**互素！！！**

![348794204-accd3992-a56e-4059-a472-39ba5ad75660](https://github.com/user-attachments/assets/a6fd2cb8-7517-4627-8fd6-0cf29711b09d)

> 点击此处了解更多关于 pad 数量选择的卡方结果：https://github.com/xtaci/qpp/blob/main/misc/chi-square.csv


正如您从 **[卡方分布](https://zh.wikipedia.org/wiki/%E5%8D%A1%E6%96%B9%E5%88%86%E5%B8%83)** 中看到的那样，通过使用与 8 互素的数字，随机性得到了增强。

## 结论

Quantum Permutation Pad 是量子密码学中一种很有前途的方法，它利用量子力学特性来实现安全通信。通过应用量子置换来加密和解密数据，QPP 在利用量子技术独特能力的同时确保了高安全性。随着量子计算和量子通信研究和技术的进步，像 QPP 这样的协议将在下一代安全通信系统中发挥至关重要的作用。

## 贡献

欢迎贡献！请提交 issue 或 pull request 以进行任何改进、错误修复或添加新功能。

## 许可证

本项目采用 GPLv3 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 参考资料

有关更详细的信息，请参阅[研究论文](https://link.springer.com/content/pdf/10.1140/epjqt/s40507-023-00164-3.pdf)。

## 致谢

特别感谢研究论文的作者在 Quantum Permutation Pad 方面的开创性工作。
