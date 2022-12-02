import { ethers } from "hardhat";

async function main() {
  const signer = (await ethers.getSigners())[0];

  console.log("signer addresse:", await signer.getAddress());
  console.log("before signer balance:", await signer.getBalance());

  const sample = await ethers.getContractFactory("Sample", signer);
  const s = await sample.deploy();

  const contract = await s.deployed();

  console.log("after signer balance:", await signer.getBalance());

  console.log(`Sample deployed to: ${s.address}`);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
