# Quantum Permutation Pad
![346540972-7a05cc4b-b207-437e-a015-57918906b738](https://github.com/user-attachments/assets/ad8ba0f1-86fa-4de8-9bbd-5fe860c40ff1)

[![GoDoc][1]][2] [![Go Report Card][3]][4]

[1]: https://godoc.org/github.com/xtaci/qpp?status.svg
[2]: https://pkg.go.dev/github.com/xtaci/qpp
[3]: https://goreportcard.com/badge/github.com/xtaci/qpp
[4]: https://goreportcard.com/report/github.com/xtaci/qpp

The [Quantum Permutation Pad](https://link.springer.com/content/pdf/10.1140/epjqt/s40507-023-00164-3.pdf) (QPP) is a cryptographic protocol designed to leverage the principles of quantum mechanics for secure communication. While the exact details of the QPP can vary based on the specific implementation and the theoretical model, the general concept involves using quantum properties such as superposition and entanglement to enhance the security of data transmission. Here’s an overview of the QPP and its relationship to quantum mechanics and cryptography:

## Key Concepts of Quantum Permutation Pad

1. **Quantum Mechanics Principles**: QPP relies on fundamental quantum mechanics principles, particularly superposition (the ability of quantum bits to be in multiple states simultaneously) and entanglement (the correlation between quantum bits regardless of distance).

2. **Quantum Bits (Qubits)**: Instead of classical bits (which are either 0 or 1), QPP uses qubits, which can be in a state of 0, 1, or any quantum superposition of these states.

3. **Permutation Operations**: Permutations in the context of QPP refer to rearranging the ways to **transform/interpret** the plaintext comparing to rearranging the keys round and round in classical cryptography. These total permutations can be thought of as $P_n$. For a 8-bit byte, the overall permutations is $P_{256} =$
```256! =857817775342842654119082271681232625157781520279485619859655650377269452553147589377440291360451408450375885342336584306157196834693696475322289288497426025679637332563368786442675207626794560187968867971521143307702077526646451464709187326100832876325702818980773671781454170250523018608495319068138257481070252817559459476987034665712738139286205234756808218860701203611083152093501947437109101726968262861606263662435022840944191408424615936000000000000000000000000000000000000000000000000000000000000000```

## Applications and Benefits

- **High Security**: QPP offers higher security levels compared to classical cryptographic methods, leveraging the unique properties of quantum mechanics.
- **Future-Proof**: As quantum computers become more powerful, classical cryptographic schemes (like RSA and ECC) are at risk. QPP provides a quantum-resistant alternative.
- **Secure Communication**: Useful for secure communications in quantum networks and for safeguarding highly sensitive data.

## Example Usage
Internal PRNG(NOT RECOMMENDED)
```golang
func main() {
    seed := make([]byte, 32)
    io.ReadFull(rand.Reader, seed)

    qpp := NewQPP(seed, 977) // a prime number of pads

    msg := make([]byte, 65536)
    io.ReadFull(rand.Reader, msg)

    qpp.Encrypt(msg)
    qpp.Decrypt(msg)
}
```

External PRNG, with shared pads. (**RECOMMENDED**)
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

## Security design in this implementatoin
The overall security is equivalent to **1683-bit** symmetric encryption.

The number of permutation matrices in an 8-qubit system is determined based on the provided seed and is selected randomly.
<img width="867" alt="image" src="https://github.com/user-attachments/assets/93ce8634-5300-47b1-ba1b-e46d9b46b432">

The permutation pad could be written in [Cycle notation](https://en.wikipedia.org/wiki/Permutation#Cycle_notation) as: $\sigma =(1\ 2\ 255)(3\ 36)(4\ 82\ 125)(....)$, which the elements are not reversible by **XORing** twice like in other [stream ciphers](https://en.wikipedia.org/wiki/Stream_cipher_attacks).

#### Locality vs. Randomness
Random permutation disrupts [locality](https://en.wikipedia.org/wiki/Principle_of_locality), which is crucial for performance. To achieve higher encryption speed, we need to maintain some level of locality. In this design, instead of switching pads for every byte, we switch to a new random pad every 8 bytes.

![349804164-3f6da444-a9f4-4d0a-b190-d59f2dca9f00](https://github.com/user-attachments/assets/2358766e-a0d3-4c21-93cb-c221aa0cece0)

The diagram clearly demonstrates that switching pads for every byte results in low performance, whereas switching pads every 8 bytes yields adequate performance.

Try directly from https://github.com/xtaci/kcptun/releases with the ```-QPP``` option enabled.

#### Performance
In modern CPUs, the latest QPP optimization can easily achieve speeds exceeding 1GB/s.
![348621244-4061d4a9-e7fa-43f5-89ef-f6ef6c00a2e7](https://github.com/user-attachments/assets/78952157-df39-4088-b423-01f45548b782)

## Security consideration in setting PADs

The number of pads should ideally be coprime with 8, as the results indicate a hidden structure in the PRNG related to the number 8.

![88d8de919445147f5d44ee059cca371](https://github.com/user-attachments/assets/9e1a160d-5433-4e24-9782-2ae88d87453d)

We demonstrate encrypting the Bible with 64 pads and 15 pads below with real data(The Bible encrypted).

**For Pads(64)**, then $GCD(64,8) == 8$

![348794146-4f6d5904-2663-46d7-870d-9fd7435df4d0](https://github.com/user-attachments/assets/e2a67bad-7d10-46e4-8d23-9866918ef04b)

**For Pads(15)**, then $GCD(15,8) == 1$, **COPRIME!!!**

![348794204-accd3992-a56e-4059-a472-39ba5ad75660](https://github.com/user-attachments/assets/a6fd2cb8-7517-4627-8fd6-0cf29711b09d)


As you can tell the difference from the **[Chi square distribution](https://en.wikipedia.org/wiki/Chi-squared_distribution)**, randomness has been enhanced by setting to numbers that are coprimes to 8.

## Conclusion

The Quantum Permutation Pad is a promising approach in the field of quantum cryptography, utilizing quantum mechanical properties to achieve secure communication. By applying quantum permutations to encrypt and decrypt data, QPP ensures high security and leverages the unique capabilities of quantum technology. As research and technology in quantum computing and quantum communication advance, protocols like QPP will play a crucial role in the next generation of secure communication systems.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements, bug fixes, or additional features.

## License

This project is licensed under the GPLv3 License. See the [LICENSE](LICENSE) file for details.

## References

For more detailed information, please refer to the [research paper](https://link.springer.com/content/pdf/10.1140/epjqt/s40507-023-00164-3.pdf).

## Acknowledgments

Special thanks to the authors of the research paper for their groundbreaking work on Quantum Permutation Pad.
