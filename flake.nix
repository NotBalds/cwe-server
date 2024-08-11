{
	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
	};

	outputs = { self, nixpkgs, ... }:
	let pkgs = nixpkgs.legacyPackages.aarch64-linux;
	xpkgs = nixpkgs.legacyPackages.x86_64-linux;
	in {
		packages.aarch64-linux.default = pkgs.buildGoModule {
			name = "cwe_server";
			src = ./.;
			vendorHash = "sha256-In/RRBbl64wBG9xLv7FnYGDOOOvYUZ52Cey6ck0mBuM=";
		};
		packages.x86_64-linux.default = xpkgs.buildGoModule {
			name = "cwe_server";
			src = ./.;
			vendorHash = "sha256-In/RRBbl64wBG9xLv7FnYGDOOOvYUZ52Cey6ck0mBuM=";
		};
		nixosModules.cwe_server = { config, lib, pkgs, ... }: {
			options = {
				server.cwe_server.enable = lib.mkEnableOption "Enable cwe server";
			};
			config = lib.mkIf config.server.cwe_server.enable {
				systemd.services.cwe_server = {
					wantedBy = [ "multi-user.target" ];
					serviceConfig = {
						WorkingDirectory = "/var/lib/cwe";
						ExecStart = "${self.packages."${pkgs.stdenv.hostPlatform.system}".default}/bin/cwe_server";
					};
				};
				services.frp.settings.proxies = [{
					name = "cwe";
					type = "tcp";
					localIP = "127.0.0.1";
					localPort = 1337;
					remotePort = 1337;
				}];
			};
		};
	};
}
