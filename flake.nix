{
	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
	};

	outputs = { self, nixpkgs, ... }:
	let system = "aarch64-linux";
	pkgs = nixpkgs.legacyPackages.${system};
	in {
		packages."${system}".default = pkgs.buildGoModule {
			name = "cwe-server";
			src = ./.;
			vendorHash = null;
		};
		nixosModules.cwe-server = { config, lib, ... }: {
			options = {
				server.cwe-server.enable = lib.mkEnableOption "Enable cwe server";
			};
			config = lib.mkIf config.server.cwe-server.enable {
				systemd.services.cwe-server = {
					wantedBy = [ "multi-user.target" ];
					serviceConfig = {
						WorkingDirectory = "${self.packages."${system}".default}";
						ExecStart = "${self.packages."${system}".default}/bin/cwe-server";
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
