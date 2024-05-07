{ src
, lib
, stdenv
, darwin
, rustPlatform
, symlinkJoin
, openssl
, rust-bin
}:
rustPlatform.buildRustPackage rec {
  inherit src;
  name = "hermes";

  nativeBuildInputs = [
    rust-bin.stable.latest.minimal
  ];

  cargoLock = {
    lockFile = "${src}/Cargo.lock";
  };

  cargoSha256 = "";
  cargoBuildFlags = "-p ibc-relayer-cli";
  buildInputs = lib.optionals stdenv.isDarwin [
    darwin.apple_sdk.frameworks.Security
    darwin.libiconv
    openssl
    darwin.libiconv
    darwin.apple_sdk.frameworks.SystemConfiguration
  ];
  doCheck = false;
  RUSTFLAGS = "--cfg ossl111 --cfg ossl110 --cfg ossl101";
  OPENSSL_NO_VENDOR = "1";
  OPENSSL_DIR = symlinkJoin {
    name = "openssl";
    paths = with openssl; [ out dev ];
  };
}
