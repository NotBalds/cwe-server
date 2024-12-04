{
	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
        flake-parts.url = "github:hercules-ci/flake-parts";
	};

    outputs = { nixpkgs, self, flake-parts, ... }@inputs: flake-parts.lib.mkFlake { inherit inputs; } {
		flake = {
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
        systems = nixpkgs.lib.platforms.linux;
        perSystem = { pkgs, ... }: {
			packages.default = pkgs.buildGoModule {
				pname = "cwe-server";
				version = "v1.0.0";
				src = ./.;
				vendorHash = null;
			};
		};
	};
}
