// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {model} from '../models';

export function DeleteGame(arg1:string):Promise<void>;

export function GetGameInfo(arg1:string,arg2:string):Promise<model.LocalGame>;

export function GetGameList(arg1:string):Promise<Array<model.GameInfo>>;

export function GetGameListPage(arg1:number):Promise<Array<model.GameInfo>>;

export function RunGame(arg1:string):Promise<void>;