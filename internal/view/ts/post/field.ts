import { PostField, StringType, StringArrayType, ImageType, ImageArrayType } from './types.js';

export class PostTitleField extends PostField {
    id: string = 'title';
    name: string = 'title';
    valueType: string = StringType;
}

export class PostContentField extends PostField {
    id: string = 'content';
    name: string = 'content';
    valueType: string = StringType;
}

export class PostTagsField extends PostField {
    id: string = 'tags[]';
    name: string = 'tags';
    valueType: string = StringArrayType;
}

export class PostPhotosField extends PostField {
    id: string = 'photos[]';
    name: string = 'files';
    valueType: string = ImageArrayType;
}