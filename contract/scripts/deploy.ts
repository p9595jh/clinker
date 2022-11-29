import { ethers } from "hardhat";

async function main() {
  const signer = (await ethers.getSigners())[0];

  console.log("signer addresse:", await signer.getAddress());
  console.log("before signer balance:", await signer.getBalance());

  const ClickerERC721 = await ethers.getContractFactory(
    "ClinkerERC721",
    signer
  );
  const clicker = await ClickerERC721.deploy();

  const contract = await clicker.deployed();

  console.log("after signer balance:", await signer.getBalance());

  console.log(`ClinkerERC721 deployed to: ${clicker.address}`);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
