import { readFileSync } from "fs";
import { createPrivateKey, createPublicKey, KeyObject } from "crypto";
import * as jose from "jose";

export type JwtSigningAlgorithm =
  | "ES256"
  | "ES384"
  | "ES512"
  | "RS256"
  | "RS384"
  | "RS512"
  | "EdDSA";

export interface JwtResources {
  privateKey: KeyObject;
  publicKey: KeyObject;
  algorithm: JwtSigningAlgorithm;
}

function detectAlgorithm(privateKey: KeyObject): JwtSigningAlgorithm {
  const keyType = privateKey.asymmetricKeyType;

  switch (keyType) {
    case "ec": {
      const keyDetails = privateKey.asymmetricKeyDetails;
      if (!keyDetails?.namedCurve) {
        throw new Error("Cannot determine EC curve");
      }
      switch (keyDetails.namedCurve) {
        case "prime256v1":
        case "P-256":
          return "ES256";
        case "secp384r1":
        case "P-384":
          return "ES384";
        case "secp521r1":
        case "P-521":
          return "ES512";
        default:
          throw new Error(`Unsupported EC curve: ${keyDetails.namedCurve}`);
      }
    }
    case "rsa": {
      const keyDetails = privateKey.asymmetricKeyDetails;
      const modulusLength = keyDetails?.modulusLength ?? 2048;
      if (modulusLength >= 4096) {
        return "RS512";
      } else if (modulusLength >= 3072) {
        return "RS384";
      }
      return "RS256";
    }
    case "ed25519":
    case "ed448":
      return "EdDSA";
    default:
      throw new Error(`Unsupported key type: ${keyType}`);
  }
}

export function loadJwtKeys(privateKeyPath: string): JwtResources {
  const privateKeyPem = readFileSync(privateKeyPath, "utf-8");
  const privateKey = createPrivateKey(privateKeyPem);
  const publicKey = createPublicKey(privateKey);
  const algorithm = detectAlgorithm(privateKey);

  return {
    privateKey,
    publicKey,
    algorithm,
  };
}

export async function signJwt(resources: JwtResources, payload: jose.JWTPayload): Promise<string> {
  const jwt = await new jose.SignJWT(payload)
    .setProtectedHeader({ alg: resources.algorithm })
    .sign(resources.privateKey);

  return jwt;
}

export async function verifyJwt(resources: JwtResources, token: string): Promise<jose.JWTPayload> {
  const { payload } = await jose.jwtVerify(token, resources.publicKey, {
    algorithms: [resources.algorithm],
  });
  return payload;
}

export async function parseJwt(
  resources: JwtResources,
  token: string
): Promise<{ payload: jose.JWTPayload; valid: boolean }> {
  try {
    const payload = await verifyJwt(resources, token);
    return { payload, valid: true };
  } catch {
    try {
      const payload = jose.decodeJwt(token);
      return { payload, valid: false };
    } catch {
      throw new Error("Invalid token format");
    }
  }
}
