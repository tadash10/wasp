// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmclient from "wasmclient"

const ArgBlobs = "this";
const ArgField = "field";
const ArgHash = "hash";

const ResBlobSizes = "this";
const ResBytes = "bytes";
const ResHash = "hash";

///////////////////////////// storeBlob /////////////////////////////

export class StoreBlobFunc extends wasmclient.ClientFunc {
	private args: wasmclient.Arguments = new wasmclient.Arguments();
	
	public blobs(v: wasmclient.Bytes): void {
		this.args.setBytes(ArgBlobs, v);
	}
	
	public async post(): Promise<wasmclient.RequestID> {
		return await super.post(0xddd4c281, this.args);
	}
}

export class StoreBlobResults extends wasmclient.ViewResults {

	hash(): wasmclient.Hash {
		return this.res.getHash(ResHash);
	}
}

///////////////////////////// getBlobField /////////////////////////////

export class GetBlobFieldView extends wasmclient.ClientView {
	private args: wasmclient.Arguments = new wasmclient.Arguments();
	
	public field(v: string): void {
		this.args.setString(ArgField, v);
	}
	
	public hash(v: wasmclient.Hash): void {
		this.args.setHash(ArgHash, v);
	}

	public async call(): Promise<GetBlobFieldResults> {
		this.args.mandatory(ArgField);
		this.args.mandatory(ArgHash);
		return new GetBlobFieldResults(await this.callView("getBlobField", this.args));
	}
}

export class GetBlobFieldResults extends wasmclient.ViewResults {

	bytes(): wasmclient.Bytes {
		return this.res.getBytes(ResBytes);
	}
}

///////////////////////////// getBlobInfo /////////////////////////////

export class GetBlobInfoView extends wasmclient.ClientView {
	private args: wasmclient.Arguments = new wasmclient.Arguments();
	
	public hash(v: wasmclient.Hash): void {
		this.args.setHash(ArgHash, v);
	}

	public async call(): Promise<GetBlobInfoResults> {
		this.args.mandatory(ArgHash);
		return new GetBlobInfoResults(await this.callView("getBlobInfo", this.args));
	}
}

export class GetBlobInfoResults extends wasmclient.ViewResults {

	blobSizes(): wasmclient.Int32 {
		return this.res.getInt32(ResBlobSizes);
	}
}

///////////////////////////// listBlobs /////////////////////////////

export class ListBlobsView extends wasmclient.ClientView {

	public async call(): Promise<ListBlobsResults> {
		return new ListBlobsResults(await this.callView("listBlobs", null));
	}
}

export class ListBlobsResults extends wasmclient.ViewResults {

	blobSizes(): wasmclient.Int32 {
		return this.res.getInt32(ResBlobSizes);
	}
}

///////////////////////////// CoreBlobService /////////////////////////////

export class CoreBlobService extends wasmclient.Service {

	public constructor(cl: wasmclient.ServiceClient) {
		super(cl, 0xfd91bc63, null);
	}

	public storeBlob(): StoreBlobFunc {
		return new StoreBlobFunc(this);
	}

	public getBlobField(): GetBlobFieldView {
		return new GetBlobFieldView(this);
	}

	public getBlobInfo(): GetBlobInfoView {
		return new GetBlobInfoView(this);
	}

	public listBlobs(): ListBlobsView {
		return new ListBlobsView(this);
	}
}
