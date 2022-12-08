import { ethers } from "hardhat";
import mongoose from "mongoose";
import Abi from "../model/abi.schema";
import Ca from "../model/ca.schema";

async function main() {
  const signer = (await ethers.getSigners())[0];

  console.log("signer addresse:", await signer.getAddress());
  console.log("before signer balance:", await signer.getBalance());

  const ClickerERC721 = await ethers.getContractFactory(
    "ClinkerERC721",
    signer
  );
  const clinker = await ClickerERC721.deploy();

  const contract = await clinker.deployed();

  console.log("after signer balance:", await signer.getBalance());
  console.log(`ClinkerERC721 deployed to: ${clinker.address}`);

  return clinker.address;
}

main()
  .then(async (address) => {
    mongoose.connect(
      `mongodb://${process.env.DB_HOST}:${process.env.DB_PORT}/${process.env.DB_SCHEMA}`
    );
    const timestamp = Math.floor(Date.now() / 1000);
    const ca = new Ca({
      address,
      timestamp,
      name: "clinker",
    });
    const { _id } = await ca.save();
    console.log("result saved into mongodb at", _id.toString());

    new Abi({
      timestamp,
      caId: _id.toString(),
      abi: require("../artifacts/contracts/ClinkerERC721.sol/ClinkerERC721.json")
        .abi,
    })
      .save()
      .then(({ _id }) =>
        console.log(`abi saved into mongodb at ${_id.toString()}`)
      )
      .finally(() => process.exit(0));
  })
  .catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
