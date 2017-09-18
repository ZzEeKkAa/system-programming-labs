<?php
include "lib/rtb_lib.php";
$redis_info = new Redis();
$redis_info->connect('localhost',9001);

$bid=60301;
$site_id='GDN';
$request_id='olo';
$user_id='uid';
$bidder = new \stdClass;
$bidder->ssp_id=3785; // comment 1

header('content-type:text/plain');

/*
 *
 * comment 2
 */

$code=get_banner_code(3785, $site_id, $request_id, $user_id, 1, $bid, $request_id,0);

print_r($code);
