import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

const config: HardhatUserConfig = {
  solidity: "0.8.17",
  networks: {
    ganache_cli: {
      url: "http://127.0.0.1:8545",
    },
    ganache_docker: {
      url: "http://127.0.0.1:9545",
    },
  },
};

export default config;
