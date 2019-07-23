/******************************************************************************
 * NTRU Cryptography Reference Source Code submitting to NIST call for
 * proposals for post quantum cryptography
 *
 * This code is written by Zhenfei Zhang @ OnboardSecurity, with additional
 * codes from public domain.
 *
 ******************************************************************************/
/*
 * api.c
 *
 *  Created on: Aug 29, 2017
 *      Author: zhenfei
 */

#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "api.h"
#include "NTRUEncrypt.h"
#include "packing.h"
#include "../common/crypto_hash_sha512.h"




/* key gen */
int crypto_encrypt_keypair(
    unsigned char       *pk,
    unsigned char       *sk)
{
    uint16_t    *F, *g, *h, *buf, *mem;
    PARAM_SET   *param;

    param   = get_param_set_by_id(TEST_PARAM_SET);

    /* memory for 3 ring elements: f, g and h */
    mem     = malloc (sizeof(uint16_t)*param->padN * 3);
    buf     = malloc (sizeof(uint16_t)*param->padN * 6);
    if (!mem )
    {
        printf("malloc error!\n");
        return -1;
    }

    F = mem;
    g = F   + param->padN;
    h = g   + param->padN;

    keygen(F,g,h,buf,param);

    /* pack h into pk */
    pack_public_key(pk, param, h);

    /* pack F,h into sk */
    pack_secret_key_CCA(sk, param, F, h);


    free(mem);
    free(buf);
    return 0;
}

/* encryption */
int crypto_encrypt(
    unsigned char       *c,
    unsigned long long  *clen,
    const unsigned char *m,
    unsigned long long  mlen,
    const unsigned char *pk)
{

    /* load the parameters */
    PARAM_SET   *param;
    uint16_t    *buf, *mem, *h, *cpoly;
    param   = get_param_set_by_id(pk[0]);

    *clen   = (unsigned long long ) param->packpk;

    if (param->id!=NTRU_CCA_443 && param->id != NTRU_CCA_743)
    {
        printf("unsupported parameter sets\n");
        return -1;
    }
    
    /* set up the memory */
    mem     = malloc(sizeof(uint16_t)*param->padN*2);
    buf     = malloc(sizeof(uint16_t)*param->padN*6);
        
    if(!mem || !buf)
    {
        printf("malloc error\n");
        return -1;
    }
    
    memset(mem,0, sizeof(uint16_t)*param->padN*2);
    memset(buf,0, sizeof(uint16_t)*param->padN*6);
    h       = mem;
    cpoly   = h     + param->padN;

    /* unpack the public key */
    unpack_public_key(pk,param, h);


    /* encryption */
    encrypt_cca(cpoly, (char*) m, mlen, h,  buf, param);

    /* pack cpoly into a ciphertext string */
    pack_public_key (c, param, cpoly);

    memset(mem,0, sizeof(uint16_t)*param->padN*2);
    memset(buf,0, sizeof(uint16_t)*param->padN*6);
    free(mem);
    free(buf);

    return 0;
}

/* decryption */
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

    if (param->id!=NTRU_CCA_443 && param->id != NTRU_CCA_743)
    {
        printf("unsupported parameter sets\n");
        return -1;
    }

    /* set up the memory */
    uint16_t    *buf, *mem, *F, *cpoly, *h;
    mem     = malloc(sizeof(uint16_t)*param->padN*4);
    buf     = malloc(sizeof(uint16_t)*param->padN*8);

    if(!mem || !buf)
    {
        printf("malloc error\n");
        return -1;
    }

    F       = mem;
    cpoly   = F     + param->padN;
    h       = cpoly + param->padN;

    memset(mem,0, sizeof(uint16_t)*param->padN*4);
    memset(buf,0, sizeof(uint16_t)*param->padN*8);


    /* unpack the keys */
    unpack_public_key (c, param, cpoly);

    unpack_secret_key_CCA (sk, param, F, h);

    /* decryption */
    *mlen = decrypt_cca((char*) m,  F, h, cpoly,  buf, param);

    free(mem);
    free(buf);

    return 0;
}

