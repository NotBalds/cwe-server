{
	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
	};

	outputs = { self, nixpkgs, ... }:
	let system = "aarch64-linux";
	pkgs = nixpkgs.legacyPackages.${system};
	in {
		packages."${system}".default = pkgs.buildGoModule {
			name = "cwe_server";
			src = ./.;
			vendorHash = "sha256-DahEqghgqBg/SL/Snu8IS8mv826otPibtseeNuHKJZU=";
		};
		packages.x86_64-linux.default = self.packages.${system}.default;
		nixosModules.cwe_server = { config, lib, ... }: {
			options = {
				server.cwe_server.enable = lib.mkEnableOption "Enable cwe server";
			};
			config = lib.mkIf config.server.cwe_server.enable {
				systemd.services.cwe_server = {
					wantedBy = [ "multi-user.target" ];
					serviceConfig = {
						WorkingDirectory = "/var/lib/cwe";
						ExecStart = "${self.packages."${system}".default}/bin/cwe_server";
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
