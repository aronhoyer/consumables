CREATE TABLE IF NOT EXISTS item_convert (
	name VARCHAR NOT NULL PRIMARY KEY,
	bom_id VARCHAR NOT NULL,
	conversion_factor REAL NOT NULL,
	default_quantity SMALLINT NOT NULL,
	destination_item_number VARCHAR NOT NULL,
	source_item_number VARCHAR NOT NULL,
	destination_unit VARCHAR NOT NULL
);
