
async function test1(){

    const keyPair1 = await createDHPair();
    const keyPair2 = await createDHPair();

    console.log("s1:", keyPair1);
    console.log("s2:", keyPair2);

    //const publicKey1 = await getPublicKey(keyPair1);
    //const publicKey2 = await getPublicKey(keyPair2);

    //console.log("public1:", publicKey1);
    //console.log("public2:", publicKey2);

    const publicKey12 = await getPublicKeyRaw(keyPair2);
    console.log("public2:", publicKey12);
    const publicKey2 = await  importPublicKey(publicKey12);
    console.log("public2:", publicKey2);


    // const privateKey1 = await getPrivateKey(keyPair1);
    // const privateKey2 = await getPrivateKey(keyPair2);
    //
    // console.log("private1:", privateKey1);
    // console.log("private2:", privateKey2);

    // 计算共享密钥
//         const sharedSecret1 = await calculateSharedSecret(keyPair1.privateKey, publicKey2);
//         const sharedSecret2 = await calculateSharedSecret(keyPair2.privateKey, keyPair1.publicKey);
//
// // sharedSecret1 和 sharedSecret2 应该是相同的（如果计算正确的话）
//         console.log("Shared secret 1:", sharedSecret1);
//         console.log("Shared secret 2:", sharedSecret2);
//
//         // 创建一个 Uint8Array
//         const uint8Array = new Uint8Array([1, 2, 3, 4, 5]);
//         saveSharedKey(sharedSecret1)
//
//         const loadKey = loadSharedKey()
//
//
// // 检查存储的数据
//         console.log("loaded=", loadKey);  // 输出：Uint8Array [1, 2, 3, 4, 5]
//
//         const keyPrint = bytesToInt64(loadKey);
//         console.log("key print is ", keyPrint)

}

