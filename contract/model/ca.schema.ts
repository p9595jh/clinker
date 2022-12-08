import { Schema, model } from "mongoose";

interface ICa {
  name: string;
  address: string;
  timestamp: number;
}

const Ca = model<ICa>(
  "Ca",
  new Schema<ICa>({
    name: { type: String, required: true },
    address: { type: String, required: true },
    timestamp: { type: Number, required: true },
  })
);

export default Ca;
