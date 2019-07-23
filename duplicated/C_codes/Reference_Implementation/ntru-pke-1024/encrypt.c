/*
 * encrypt.c
 *
 *  Created on: Aug 31, 2017
 *      Author: zhenfei
 */



#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "api.h"
#include "NTRUEncrypt.h"
#include "../common/crypto_hash_sha512.h"



/* key gen */
int crypto_encrypt_keypair(
    unsigned char       *pk,
    unsigned char       *sk)
{
    int64_t     *f, *g, *hntt, *buf, *mem;
    PARAM_SET   *param;
    param   = get_param_set_by_id(TEST_PARAM_SET);

    /* memory for 3 ring elements: f, g and h */
    mem     = malloc (sizeof(int64_t)*param->N * 3);
    buf     = malloc (sizeof(int64_t)*param->N * 2);
    if (!mem || !buf)
    {
        printf("malloc error!\n");
        return -1;
    }

    f       = mem;
    g       = f   + param->N;
    hntt    = g   + param->N;

    keygen(f,g,hntt,buf,param);

    /* pack h into pk */
    pack_ring_element(pk, param, hntt);

    /* pack F into sk */
    pack_ring_element(sk, param, f);
    pack_ring_element(sk+param->N*sizeof(int32_t)/sizeof(unsigned char)+1, param, hntt);

    memset(mem,0, sizeof(int64_t)*param->N*3);
    memset(buf,0, sizeof(int64_t)*param->N*2);

    free(mem);
    free(buf);

    return 0;
}

/* ebacs API: encryption */
int crypto_encrypt(
    unsigned char       *c,
    unsigned long long  *clen,
    const unsigned char *m,
    unsigned long long  mlen,
    const unsigned char *pk)
{
    /* load the parameters */
    PARAM_SET   *param;
    param   = get_param_set_by_id(pk[0]);

    if (param->id!=NTRU_CCA_1024)
    {
        printf("unsupported parameter sets\n");
        return -1;
    }
    
    /* set up the memory */

    int64_t    *buf, *mem, *hntt, *cpoly;
    mem     = malloc(sizeof(int64_t)*param->N*2);
    buf     = malloc(sizeof(int64_t)*param->N*7 + LENGTH_OF_HASH*2);
    if (!mem || !buf )
    {
        printf("malloc error!\n");
        return -1;
    }

    hntt    = mem;
    cpoly   = hntt  + param->N;


    memset(mem,0, sizeof(int64_t)*param->N*2);
    memset(buf,0, sizeof(int64_t)*param->N*7 + LENGTH_OF_HASH*2);


    /* unpack the public key */
    unpack_ring_element(pk, param, hntt);

    /* encryption */
    encrypt_cca(cpoly, (char*) m, mlen, hntt, buf, param);


    /* pack cpoly into a ciphertext string */
    pack_ring_element (c, param, cpoly);

    *clen = param->N*sizeof(int32_t)/sizeof(unsigned char)+1;


    memset(mem,0, sizeof(int64_t)*param->N*2);
    memset(buf,0, sizeof(int64_t)*param->N*7 + LENGTH_OF_HASH*2);
    free(mem);
    free(buf);


    return 0;
}

/* ebacs API: decryption */
int crypto_encrypt_open(
    unsigned char       *m,
    unsigned long long  *mlen,
    const unsigned char *c,
    unsigned long long  clen,
    const unsigned char *sk)
{
    /* load the parameters */
    PARAM_SET   *param;

    param   =   get_param_set_by_id(c[0]);
    if (param->id!=NTRU_CCA_1024)
    {
        printf("unsupported parameter sets\n");
        return -1;
    }

    /* set up the memory */
    int64_t    *buf, *mem, *f, *cpoly, *mpoly, *hntt;
    mem     = malloc(sizeof(int64_t)*param->N*4);
    buf     = malloc(sizeof(int64_t)*param->N*7 + LENGTH_OF_HASH*2);

    if (!mem || !buf )
    {
        printf("malloc error!\n");
        return -1;
    }

    f       = mem;
    cpoly   = f     + param->N;
    mpoly   = cpoly + param->N;
    hntt    = mpoly + param->N;

    memset(mem,0, sizeof(int64_t)*param->N*4);
    memset(buf,0, sizeof(int64_t)*param->N*7 + LENGTH_OF_HASH*2);

    /* unpack the keys */
    unpack_ring_element (c, param, cpoly);

    unpack_ring_element (sk, param, f);

    unpack_ring_element (sk+param->N*sizeof(int32_t)/sizeof(unsigned char)+1, param, hntt);

    /* decryption */
    *mlen = decrypt_cca((char*) m, f, hntt, cpoly, buf, param);


    memset(mem,0, sizeof(int64_t)*param->N*3);
    memset(buf,0, sizeof(int64_t)*param->N*7 + LENGTH_OF_HASH*2);
    free(mem);
    free(buf);


    return 0;
}

