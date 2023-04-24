public void serialize(LittleEndianOutput out) {out.writeShort(field_1_vcenter);}
public void addAll(BlockList<T> src) {if (src.size == 0)return;int srcDirIdx = 0;for (; srcDirIdx < src.tailDirIdx; srcDirIdx++)addAll(src.directory[srcDirIdx], 0, BLOCK_SIZE);if (src.tailBlkIdx != 0)addAll(src.tailBlock, 0, src.tailBlkIdx);}
public void writeByte(byte b) {if (upto == blockSize) {if (currentBlock != null) {addBlock(currentBlock);}currentBlock = new byte[blockSize];upto = 0;}currentBlock[upto++] = b;}
public ObjectId getObjectId() {return objectId;}
public DeleteDomainEntryResult deleteDomainEntry(DeleteDomainEntryRequest request) {request = beforeClientExecution(request);return executeDeleteDomainEntry(request);}
public long ramBytesUsed() {return ((termOffsets!=null)? termOffsets.ramBytesUsed() : 0) +((termsDictOffsets!=null)? termsDictOffsets.ramBytesUsed() : 0);}
public final String getFullMessage() {byte[] raw = buffer;int msgB = RawParseUtils.tagMessage(raw, 0);if (msgB < 0) {return ""; }return RawParseUtils.decode(guessEncoding(), raw, msgB, raw.length);}
public POIFSFileSystem() {this(true);_header.setBATCount(1);_header.setBATArray(new int[]{1});BATBlock bb = BATBlock.createEmptyBATBlock(bigBlockSize, false);bb.setOurBlockIndex(1);_bat_blocks.add(bb);setNextBlock(0, POIFSConstants.END_OF_CHAIN);setNextBlock(1, POIFSConstants.FAT_SECTOR_BLOCK);_property_table.setStartBlock(0);}
public void init(int address) {slice = pool.buffers[address >> ByteBlockPool.BYTE_BLOCK_SHIFT];assert slice != null;upto = address & ByteBlockPool.BYTE_BLOCK_MASK;offset0 = address;assert upto < slice.length;}
public SubmoduleAddCommand setPath(String path) {this.path = path;return this;}
