<?php
	error_reporting(E_ALL^E_NOTICE);
	global $t0;

	$t0=microtime(1);
	
	// Для дебага
	//error_log("Time: ".(round((microtime(1) - $t0), 3))." ms");
	if (0)
	{
		// IF FUCKUP (1)

		set_no_banner(); die();

	};

	$config_dir = "cfg/";

	$z_debug = false;
	if(isset($_GET["z_debug"])){
		$z_debug = true;
	}
/**
 * @todo: константа для пути инклуда - должна еще ускорить все инклуды
 */
/*	if (!defined('PATH'))
		define('PATH', __DIR__);*/

	$status_file = sys_get_temp_dir().'/rtb_server.txt';
	$status_file=str_replace("//","/", $status_file);

	if (0 && !file_exists ($status_file))
	{
		set_no_banner(null, "RTB Service is offline"); die();
	};

	if (0 && "keep-alive"!=strtolower($_SERVER["HTTP_CONNECTION"]))
	{
		set_no_banner(null, "No keep-alive detected"); die();
	};

	// we need to respond to all smaato auctions to pass the tests. For others block responses under heavy load.
	//$load = sys_getloadavg();
	//{
	//	// system overloaded
	//
	//	set_no_banner(null, "System overloaded"); die();
	//	save_log_txt('NO');
	//
	//	die();
	//};

	if(!$redis_response->exists("response:".$response_id)){
		set_no_banner($request, "Log response saving error");
	}
	
	
	header('X-bid-resp: 1');
	$bidder->incrSelfAudit('call:bid_request');
	$bidder->incrPlannerStat([], [1]);
	
	$t2=microtime(1);

	header('Cache-Control: no-cache, must-revalidate');
	header('Expires: Mon, 26 Jul 1997 05:00:00 GMT');
	

	/**
	 * @todo: перенести в универсальный биддер и биддер для гугла
	 */


	if ("77.91.130.160"!=$_SERVER['REMOTE_ADDR'])
	{
		
		if ($bidder->is_google)
		{
			$response_encode = response_json2pb(json_encode($google_code));
			header('Content-Type: application/octet-stream');
			
		}
		elseif (!isset($ar_code))
		{
			$response_encode = json_encode($response);
			header ('Content-Type: application/json; charset=UTF-8');

		}
		else
		{
			$response_encode = json_encode($ar_code);
			header ('Content-Type: application/json; charset=UTF-8');
		};

		$size=strlen($response_encode);

//		// zeka's debug
//		if($z_debug){
//			error_log($response_encode);
//		}
		echo $response_encode;
	}
	else
	{
		//echo $t2-$t1;
		//$response['ext']['tt']=array('test'=>'teststring here');
		
		$response['seatbid'][0]['bid'][0]['price_uah']=$cpm_uah;
		$response['microtime']=$t2-$t1;

		if ($bidder->is_google)
		{
			$response_encode = response_json2pb(json_encode($google_code));
			//$response_encode = (json_encode($google_code));
		}
		elseif (!isset($ar_code))
		{
			$response_encode = json_encode($response);
		}
		else
		{
			$response_encode = json_encode($ar_code);
		};
		;
		$size=strlen($response_encode);
		header("X-mt: ").$response['microtime'];
		header("Content-Length: $size");

		echo $response_encode;

		//die();
	};

	//error_log("Time before FCGI: ".(round((microtime(1) - $t0), 3))." ms");
	

	$work = 'worker_mongolog'; 
	$_msg = array();
	$_msg[$work]['_id'] = $bidder->response_id;
	$_msg[$work]['mongoCollection'] = 'dspRes';
	$_msg[$work]['sys']['ssp_id'] = $ssp_id;
	$_msg[$work]['sys']['ssp_request_id'] = isset($request['imp']['0']['id']) ? $request['imp']['0']['id'] : null;
	$_msg[$work]['sys']['c8dmp'] = $request_dmp;
	$_msg[$work]['sys']['ts'] = time();
	$_msg[$work]['request'] = $response;

	function add_data_to_bannerbid_mongo_inc($banner_id, $ssp_id, $site_id, $first=false){
		global $request;
		
		try{

			require_once("/var/www/traits/C8SimpleDBConnector.php");
			$redis_mongo_inc = C8SimpleDBConnector::connectToRedis("redis_mongo_inc");
			
			if(!$first){
			$key="bannerbid:inc:".$banner_id;
				
			$redis_mongo_inc->sAdd('mongo_bannerbid_keys', $key);
			}else{
				$key="bannerbidfirst:inc:".$banner_id;
				
				$redis_mongo_inc->sAdd('mongo_bannerbidfirst_keys', $key);
			}
			$redis_mongo_inc->hIncrBy($key, $ssp_id."_".$site_id,1);
			
			if($ssp_id==3785){

				/* OLD */
				//$redis_request_to_bids=C8RedisRegistry::getInstance()->getRedis("redis_request_to_bids");

				$redis_request_to_bids = C8SimpleDBConnector::connectToRedis("redis_request_to_bids");
				$redis_request_to_bids->set($request["id"], $banner_id);
				$redis_request_to_bids->set($request["id"]."_site", $site_id);
				$redis_request_to_bids->expire($request["id"], 1200);
				$redis_request_to_bids->expire($request["id"]."_site", 1200);
			}
		
		}catch (Exception $e){
			error_log($e->getMessage());
			return false;
		}
		
	}
	
/*
user_id

Sets (списки баннеров):

cpm_net_1_site_2 U (cpm_net_1_theme_2\cpm_net_1_site_2) U (cpm_net_1\cpm_net_1_site_2\cpm_net_1_theme_2) U (cpm\cpm_net_1\cpm_net_1_site_2\cpm_net_1_theme_2)
*/
?>
