import { model, Schema } from "mongoose";

interface IAbi {
  caId: string;
  abi: object;
  timestamp: number;
}

const Abi = model<IAbi>(
  "Abi",
  new Schema<IAbi>({
    caId: { type: String, required: true },
    abi: { type: Object, required: true },
    timestamp: { type: Number, required: true },
  })
);

export default Abi;
