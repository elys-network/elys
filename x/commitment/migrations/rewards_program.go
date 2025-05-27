package migrations

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/v5/x/commitment/types"
)

// Additional addresses for testnet only
var RewardProgramTestnet = []types.RewardProgram{
	{
		Address: "elys17lgqaj8z7hxewc2wncrf4as33l7hdja7ttn4c6",
		Amount:  math.NewInt(1000000000),
		Claimed: false,
	},
	{
		Address: "elys1u8c28343vvhwgwhf29w6hlcz73hvq7lwxmrl46",
		Amount:  math.NewInt(1000000000),
		Claimed: false,
	},
	{
		Address: "elys130yaavkws04nt9lslqc5lgyj3jpu8t20x63pcj",
		Amount:  math.NewInt(1944000000),
		Claimed: false,
	},
}

var RewardProgram = []types.RewardProgram{
	{
		Address: "elys13w62vn9h3eeecampw5sq9zjtpcwmwmjwzwjd9e",
		Amount:  math.NewInt(82639700000),
		Claimed: false,
	},
	{
		Address: "elys1jq3rknpvfmgedsynq4dad8akfmvcsyw8udfvt2",
		Amount:  math.NewInt(82087100000),
		Claimed: false,
	},
	{
		Address: "elys123scuwjkz9v778sruaymudvvsfgv0a8wdepk74",
		Amount:  math.NewInt(18148800000),
		Claimed: false,
	},
	{
		Address: "elys1tnk45uvxvdsemhd0pnl5y72jl43lwwm0wp0sls",
		Amount:  math.NewInt(17321900000),
		Claimed: false,
	},
	{
		Address: "elys174vz4a2rxme7j04g34pyxvqnu05dcsupdjxz70",
		Amount:  math.NewInt(14760900000),
		Claimed: false,
	},
	{
		Address: "elys1vlldhafw0mqftljsd2rmqgv4zww26n72ezmj39",
		Amount:  math.NewInt(12590900000),
		Claimed: false,
	},
	{
		Address: "elys1egdrqseqdhzk67ssfc7y7d72ka00w5g44a6ha7",
		Amount:  math.NewInt(12030800000),
		Claimed: false,
	},
	{
		Address: "elys1hyapd7e50vrtkgujd2h9pqa4l58r3uj524zugt",
		Amount:  math.NewInt(9114200000),
		Claimed: false,
	},
	{
		Address: "elys1pwtf4vj6mghsad3mwwjk8tjpjzlu54cdl54hjk",
		Amount:  math.NewInt(7546600000),
		Claimed: false,
	},
	{
		Address: "elys1u5gaduhhs0ht3aqj4u74wfpxrqn0nwuq496vzl",
		Amount:  math.NewInt(6797700000),
		Claimed: false,
	},
	{
		Address: "elys18ym9ke5xpvvhgtrzqu8h5p365cfyc6w8ahaqzr",
		Amount:  math.NewInt(6274000000),
		Claimed: false,
	},
	{
		Address: "elys12hwjyrdp8mzk54ut6p4teslq3ksslwrpyvhazq",
		Amount:  math.NewInt(5952400000),
		Claimed: false,
	},
	{
		Address: "elys1a7yx5k57ee7vpa3ut7gus66ak26q3mgcmx0z62",
		Amount:  math.NewInt(5077700000),
		Claimed: false,
	},
	{
		Address: "elys1t7cd7psexmemxurzgpk90pwz78vj7ffk2px7uv",
		Amount:  math.NewInt(4969500000),
		Claimed: false,
	},
	{
		Address: "elys1pg2l3len5ueu2anchzchwnlm2msus4u2ng2ayh",
		Amount:  math.NewInt(4800700000),
		Claimed: false,
	},
	{
		Address: "elys190k0y8m6pktew6xs58q2dgfarqyxj3zsahpkpg",
		Amount:  math.NewInt(4523300000),
		Claimed: false,
	},
	{
		Address: "elys1gwfuspdtuy3dat67d4yw5n7vnfytkr5agphxpw",
		Amount:  math.NewInt(4384200000),
		Claimed: false,
	},
	{
		Address: "elys1clqk8rd6vrahvlw9mxppqjj6n9gqa2mdwha3v3",
		Amount:  math.NewInt(4336000000),
		Claimed: false,
	},
	{
		Address: "elys1cvwx28t7v5s34mwyf47slu8cwcmhw0zf5fzx5y",
		Amount:  math.NewInt(4258300000),
		Claimed: false,
	},
	{
		Address: "elys1pplh4m6jpavtl0q49x2v7dw2mthvc7zu6fp6ru",
		Amount:  math.NewInt(4217600000),
		Claimed: false,
	},
	{
		Address: "elys1ada2rnxya25f84zfz6s4wrn26kjqaa6kv0w5dm",
		Amount:  math.NewInt(4177100000),
		Claimed: false,
	},
	{
		Address: "elys1pp8pxjulhetpf0ewt3w0735v7zzeg6z8vs7hr0",
		Amount:  math.NewInt(4105600000),
		Claimed: false,
	},
	{
		Address: "elys1sxyvk4eqd9m8rrn2qst67waycj3mkka2fkpjfp",
		Amount:  math.NewInt(4041200000),
		Claimed: false,
	},
	{
		Address: "elys1tcfusja8ekyxjkmguvfyfm5jl2jl20pxmyznaj",
		Amount:  math.NewInt(3943200000),
		Claimed: false,
	},
	{
		Address: "elys1pea007hmt7scdwcduecqhsr5lws72595dgpsnr",
		Amount:  math.NewInt(3827800000),
		Claimed: false,
	},
	{
		Address: "elys1ev3kufgy264nedktpw0w02smr2z9svc3dk0ucq",
		Amount:  math.NewInt(3745900000),
		Claimed: false,
	},
	{
		Address: "elys19nw9u3hdldj0s93prcvk6pz7tsv4fae54a2atp",
		Amount:  math.NewInt(3658700000),
		Claimed: false,
	},
	{
		Address: "elys1j23m9zw2xyhk4utkd32ve3xl9up4xu7lq75ple",
		Amount:  math.NewInt(3557900000),
		Claimed: false,
	},
	{
		Address: "elys1vm77qyclxzd84acgreg7nsl26jtsqwzqenwj2m",
		Amount:  math.NewInt(3484200000),
		Claimed: false,
	},
	{
		Address: "elys1v7xyjj2we4udmyyy68hhlv5gtamcut7cxnhz3t",
		Amount:  math.NewInt(3404500000),
		Claimed: false,
	},
	{
		Address: "elys1u8c28343vvhwgwhf29w6hlcz73hvq7lwxmrl46",
		Amount:  math.NewInt(3271100000),
		Claimed: false,
	},
	{
		Address: "elys1jsvqm8acqrxjhpxvuttq5wq2w6tzqtfzm73j2u",
		Amount:  math.NewInt(3161200000),
		Claimed: false,
	},
	{
		Address: "elys1dauy4vfjl2z3a4pf264luuyjkssd23980gr74t",
		Amount:  math.NewInt(3148200000),
		Claimed: false,
	},
	{
		Address: "elys1pw9tlvfk5eatz3hajp8sd5ekd9g57g0pm4eycy",
		Amount:  math.NewInt(3130900000),
		Claimed: false,
	},
	{
		Address: "elys1538jcdcc4kxwx53w4k7avpsrr38300clnk8rvn",
		Amount:  math.NewInt(3115000000),
		Claimed: false,
	},
	{
		Address: "elys1lckk5h87crvyzsc3au8sw7q0fuv9mk5n3hxwsk",
		Amount:  math.NewInt(3105900000),
		Claimed: false,
	},
	{
		Address: "elys1v0ulf080nvcfu36ln8t6rkr4mg45mks35e6smm",
		Amount:  math.NewInt(3096700000),
		Claimed: false,
	},
	{
		Address: "elys106ps2h8y6mpnvl4f5lpjna87pv0ajasasllk96",
		Amount:  math.NewInt(3015800000),
		Claimed: false,
	},
	{
		Address: "elys1zwqv4pk2smtqs5s40r5z34f9v239zu2egl9put",
		Amount:  math.NewInt(2937800000),
		Claimed: false,
	},
	{
		Address: "elys1zulfe6nr8w5lgw432pv663x8jx3uzwu6nkgllv",
		Amount:  math.NewInt(2875000000),
		Claimed: false,
	},
	{
		Address: "elys1y7dkgrprznn2pn9nmgm6evues5luncmn8z6sje",
		Amount:  math.NewInt(2720100000),
		Claimed: false,
	},
	{
		Address: "elys1v05f642ryu39020gmgfgelqadvr48uygdsvz42",
		Amount:  math.NewInt(2606500000),
		Claimed: false,
	},
	{
		Address: "elys1f398n49k6qtphmh7hccs7rxaryng78z8hxsj6g",
		Amount:  math.NewInt(2583600000),
		Claimed: false,
	},
	{
		Address: "elys1z4eat2vjfherr8eq6rzxgjvauswgyytc8qqky8",
		Amount:  math.NewInt(2566800000),
		Claimed: false,
	},
	{
		Address: "elys1qh6mh703sfpv0dggychrtfxjuvfpq8rltuq0wn",
		Amount:  math.NewInt(2500500000),
		Claimed: false,
	},
	{
		Address: "elys14r6f55m97j7gaqk4k9lt7cgmy59284kea3d6eu",
		Amount:  math.NewInt(2500000000),
		Claimed: false,
	},
	{
		Address: "elys14tqz4zr4f8v8jwct93ard76vupdphe9q34cw46",
		Amount:  math.NewInt(2500000000),
		Claimed: false,
	},
	{
		Address: "elys1cl6whjcch5mfjwnty6ltgawrdz0fvhca589klr",
		Amount:  math.NewInt(2500000000),
		Claimed: false,
	},
	{
		Address: "elys1k778ncy9hl0ssanjg80sxlugvkk7a9j0v4ew6m",
		Amount:  math.NewInt(2500000000),
		Claimed: false,
	},
	{
		Address: "elys1s84s8z9qd93uja646m3c8j53my8ykt3j64gjxa",
		Amount:  math.NewInt(2469900000),
		Claimed: false,
	},
	{
		Address: "elys1kxgan4uq0m8gqztd09n6qm627p6v4ayzngxyx8",
		Amount:  math.NewInt(2328600000),
		Claimed: false,
	},
	{
		Address: "elys1we50lf9sw27w6d33z5mukygdcldyganh8w29kh",
		Amount:  math.NewInt(2327500000),
		Claimed: false,
	},
	{
		Address: "elys1ex7nd785vcg3nspssrpvkqfww4yht4q9ulclhw",
		Amount:  math.NewInt(2302400000),
		Claimed: false,
	},
	{
		Address: "elys1st2f0myepljrv84devj4zmy9585mel6p93c46y",
		Amount:  math.NewInt(2185800000),
		Claimed: false,
	},
	{
		Address: "elys1vm8l8pan8wy9xfvfkv0aqltskp35gljqfymm7z",
		Amount:  math.NewInt(2133800000),
		Claimed: false,
	},
	{
		Address: "elys16jj50pk2ejxewrt29c0u7m3cyetepged3z55k0",
		Amount:  math.NewInt(2109199999),
		Claimed: false,
	},
	{
		Address: "elys1lmhqp9y64rdvkwqqrjx4wr8jq8hj6wz2mz0due",
		Amount:  math.NewInt(2011900000),
		Claimed: false,
	},
	{
		Address: "elys19tcd54wjxtprz4zrh30sh6drg0f4sk8dcqf6rs",
		Amount:  math.NewInt(1944000000),
		Claimed: false,
	},
	{
		Address: "elys1nfrwgu805mvwcfecvrksvrrld44yawtrmdwyqg",
		Amount:  math.NewInt(1939000000),
		Claimed: false,
	},
	{
		Address: "elys18lclxxjghhfum9tj5xyy3gg3yguypj8xve6r7g",
		Amount:  math.NewInt(1908900000),
		Claimed: false,
	},
	{
		Address: "elys1seqghatqlqakest4uhcw9728dfddke8hpp802u",
		Amount:  math.NewInt(1908800000),
		Claimed: false,
	},
	{
		Address: "elys1dtj2cpgkm77xafen6ty2hxdrqlcc8knhndvd0v",
		Amount:  math.NewInt(1895200000),
		Claimed: false,
	},
	{
		Address: "elys17mtanegq5qj3kg48ff9zcvlt2prgu3ttpghmgs",
		Amount:  math.NewInt(1875000000),
		Claimed: false,
	},
	{
		Address: "elys1gpaxtud6ch0t7q7lw08gd6ujmjrhlptejywgy3",
		Amount:  math.NewInt(1875000000),
		Claimed: false,
	},
	{
		Address: "elys1tvfggynym5nmh8rldvpch90g4v6x7dzdcrs5ea",
		Amount:  math.NewInt(1875000000),
		Claimed: false,
	},
	{
		Address: "elys174tvh2dty7vsvwn2cfsmkwq8tplqgr5fj77dds",
		Amount:  math.NewInt(1868000000),
		Claimed: false,
	},
	{
		Address: "elys1jtjny4dnclqsrx6fzcx5dka2cng3zdwk5tyrmz",
		Amount:  math.NewInt(1848400000),
		Claimed: false,
	},
	{
		Address: "elys1f2t5nggzlgzuyttj52xdtvetfqswfhfvqz2ea8",
		Amount:  math.NewInt(1820300000),
		Claimed: false,
	},
	{
		Address: "elys1aqrm0jf833qaqnf2fepd5n38tfwy549xzhptwc",
		Amount:  math.NewInt(1626100000),
		Claimed: false,
	},
	{
		Address: "elys1gjw8tnjgr6cqkatcrjz72k380f5vqhk4z3h3h3",
		Amount:  math.NewInt(1618400000),
		Claimed: false,
	},
	{
		Address: "elys15qn5eg2f6mhghqnlva89txrx0hee9gnmlcjt50",
		Amount:  math.NewInt(1564600000),
		Claimed: false,
	},
	{
		Address: "elys1s3txxl96g4y5dap7pllvtw8v2whpfs9vxxdaqq",
		Amount:  math.NewInt(1563500000),
		Claimed: false,
	},
	{
		Address: "elys17s34wcwss4hwkh4dgn95cvwlzgdpcf7msrtkjn",
		Amount:  math.NewInt(1556900000),
		Claimed: false,
	},
	{
		Address: "elys1k85qxq4pltlu5udtk8erhhmm8stelvz6ak2faq",
		Amount:  math.NewInt(1535900000),
		Claimed: false,
	},
	{
		Address: "elys1qa5fzmvrfptplzkert58pgzxjtllmd2cqtv5am",
		Amount:  math.NewInt(1535500000),
		Claimed: false,
	},
	{
		Address: "elys1z7ml3yp2sqguptlehv2fmhexpjw8sqsvcjuaa0",
		Amount:  math.NewInt(1510200000),
		Claimed: false,
	},
	{
		Address: "elys1dqn7t9dqtgsjw3m9zvnfkzf6pspglc07shkt3g",
		Amount:  math.NewInt(1501700000),
		Claimed: false,
	},
	{
		Address: "elys1nfuyz0656xfpfqq38nrnly6gxjg5nhwg9cye7g",
		Amount:  math.NewInt(1500000000),
		Claimed: false,
	},
	{
		Address: "elys1et9km7ayjskg7teq437rft8zu3s45gs88mx8ut",
		Amount:  math.NewInt(1496000000),
		Claimed: false,
	},
	{
		Address: "elys180j36lh2axv6763xg0gd334yr6au727tkw8c8k",
		Amount:  math.NewInt(1494300000),
		Claimed: false,
	},
	{
		Address: "elys1skhzdm3n4h9e698hsgv6l67luncw354fu52zyj",
		Amount:  math.NewInt(1472200000),
		Claimed: false,
	},
	{
		Address: "elys1v8h45qelm6r7ucsdmex48ntkvnrnq9urjm3w2j",
		Amount:  math.NewInt(1467100000),
		Claimed: false,
	},
	{
		Address: "elys1ggftlp28yrw4rj3ur4endwsyusua46tln9uyv7",
		Amount:  math.NewInt(1466000000),
		Claimed: false,
	},
	{
		Address: "elys15lsrfjry5f4peqwjza805pj0ucgped2ta79gwc",
		Amount:  math.NewInt(1462400000),
		Claimed: false,
	},
	{
		Address: "elys12y0a0nttq3a2vp3kf0hvt2xw0kd8entkyza329",
		Amount:  math.NewInt(1443300000),
		Claimed: false,
	},
	{
		Address: "elys14tdlfefe9w9nuhmxptj6h37gule9ndxs3y9py4",
		Amount:  math.NewInt(1441400000),
		Claimed: false,
	},
	{
		Address: "elys1y40yywg94e6hw80464kxrr7063d35ekyd23k2p",
		Amount:  math.NewInt(1411000000),
		Claimed: false,
	},
	{
		Address: "elys10dv5h692pzdjf0z3zel2hant5axqhcjugn0fu6",
		Amount:  math.NewInt(1378300000),
		Claimed: false,
	},
	{
		Address: "elys15n8ugf24lhhtrtt0up5lfhy8kezf084mnuayez",
		Amount:  math.NewInt(1376200000),
		Claimed: false,
	},
	{
		Address: "elys1rv5llx7kxd0endqeqtspju7x9mx6m9wtzclmxc",
		Amount:  math.NewInt(1360500000),
		Claimed: false,
	},
	{
		Address: "elys1wzele4enmzppy8pd6rxltsp7c0pjftqn7ml6at",
		Amount:  math.NewInt(1357800000),
		Claimed: false,
	},
	{
		Address: "elys1ukkkpp6gd2qtw4qtar5m4e2vm462p3uf2fk04e",
		Amount:  math.NewInt(1304400000),
		Claimed: false,
	},
	{
		Address: "elys1wygsp6depahahyg8zn0q5w358gceh5sy7wf6dw",
		Amount:  math.NewInt(1284700000),
		Claimed: false,
	},
	{
		Address: "elys1g5sedsmcpg62z8ce6jtktm5nmfjnmdkt9mrcf2",
		Amount:  math.NewInt(1273700000),
		Claimed: false,
	},
	{
		Address: "elys1yvsp4fvj6wwh9wwdqrj4l86nf2sz4czgdxv4yk",
		Amount:  math.NewInt(1254700000),
		Claimed: false,
	},
	{
		Address: "elys1l4cz20qf8plzva6p0g3q8q84slmn5gefqr9wrc",
		Amount:  math.NewInt(1213700000),
		Claimed: false,
	},
	{
		Address: "elys1ntkl7a7r5hcvfty7p4k4f5ep8yh3d02lsgu9y9",
		Amount:  math.NewInt(1213600000),
		Claimed: false,
	},
	{
		Address: "elys1rvtlre96vmf57acak86yvfh685gahp06pqtjxf",
		Amount:  math.NewInt(1197700000),
		Claimed: false,
	},
	{
		Address: "elys1c69yrjgrka9yafg7wcx4a0kwru9jwvwq8kf7ek",
		Amount:  math.NewInt(1189100000),
		Claimed: false,
	},
	{
		Address: "elys125nxgp2um4npqe639qepgtartutg5ykuhjaalu",
		Amount:  math.NewInt(1179700000),
		Claimed: false,
	},
	{
		Address: "elys1h6r7sgwxfxps4payfyc8rl56svzmx6t5kp9u9n",
		Amount:  math.NewInt(1166200000),
		Claimed: false,
	},
	{
		Address: "elys1gsj0jh9z9mynemm4vjyc88wldw2trtg68hs7q3",
		Amount:  math.NewInt(1164500000),
		Claimed: false,
	},
	{
		Address: "elys1nstj4l63gjl4gg4kuar77t0p4vdk48q3tzq504",
		Amount:  math.NewInt(1161900000),
		Claimed: false,
	},
	{
		Address: "elys1pudghjgk2c3sjncw2mygwdd26r6jwv20eqd4f5",
		Amount:  math.NewInt(1160500000),
		Claimed: false,
	},
	{
		Address: "elys1ua08kamwz74ue6ytcz09nrrs5nuk5jf35ge88n",
		Amount:  math.NewInt(1157200000),
		Claimed: false,
	},
	{
		Address: "elys1xslq2hmp3s3sxeggl8a56ga5jve06mttpju5sm",
		Amount:  math.NewInt(1141600000),
		Claimed: false,
	},
	{
		Address: "elys1l7dqvqfwux7zkrkxtarp3nvk0u9cz7pcftevxz",
		Amount:  math.NewInt(1132100000),
		Claimed: false,
	},
	{
		Address: "elys1zwtv7pplzh7m3t4m2zel8c0q78qsyptv3t60uy",
		Amount:  math.NewInt(1125900000),
		Claimed: false,
	},
	{
		Address: "elys1a6vmp4wyj5tads7u7x9v90e25cmmczddl9k0sn",
		Amount:  math.NewInt(1124200000),
		Claimed: false,
	},
	{
		Address: "elys1r3vr9rg7kpamyk3dwt4qmkgnr8ygh4htge7tcp",
		Amount:  math.NewInt(1121400000),
		Claimed: false,
	},
	{
		Address: "elys1wp9w68jarn9pnsfc8sxnyqpe3sue20rhj4c3fq",
		Amount:  math.NewInt(1119800000),
		Claimed: false,
	},
	{
		Address: "elys1m3myrk7tvfg89af39ppdft3de2qek7c90ra632",
		Amount:  math.NewInt(1110600000),
		Claimed: false,
	},
	{
		Address: "elys1n0t2sj4uy2ezcwswqn35eqmtcn2w0dkjlxe926",
		Amount:  math.NewInt(1108500000),
		Claimed: false,
	},
	{
		Address: "elys17md6vnlpulq2vt8u5rccadx9wxx0al7q9mkcz8",
		Amount:  math.NewInt(1081200000),
		Claimed: false,
	},
	{
		Address: "elys1hrghtnw9j9zs9ldy9cfme3ajkrcmqe7dz5vzq7",
		Amount:  math.NewInt(1079700000),
		Claimed: false,
	},
	{
		Address: "elys1uxnjvt5myz5ahrjtwkzms8ndk66s0arrmdatcy",
		Amount:  math.NewInt(1077900000),
		Claimed: false,
	},
	{
		Address: "elys16ejzx2475lz3yw824h4sf2j8frz4tq0hvehwfw",
		Amount:  math.NewInt(1064200000),
		Claimed: false,
	},
	{
		Address: "elys10aeh6eshk4rshknrgxy7uy75pdekjcvta65f0c",
		Amount:  math.NewInt(1062700000),
		Claimed: false,
	},
	{
		Address: "elys1tdxu5dvg9pq5dj0lgndcushwg8yx7xmkrm8sqf",
		Amount:  math.NewInt(1062700000),
		Claimed: false,
	},
	{
		Address: "elys14kx92jz0dpuh4gakxzr6zm97djdvjxq5tcj4t7",
		Amount:  math.NewInt(1062300000),
		Claimed: false,
	},
	{
		Address: "elys1aspx97m58jxyw42f25q2vygz50yk26c0gcgftt",
		Amount:  math.NewInt(1046900000),
		Claimed: false,
	},
	{
		Address: "elys1lgk09phf7qyvs6lkl56madxmkm7m363e7fd5fk",
		Amount:  math.NewInt(1041300000),
		Claimed: false,
	},
	{
		Address: "elys1zp7c4gjm46z9jwrdm8zlsdu993slh4rvmgpu96",
		Amount:  math.NewInt(1040599999),
		Claimed: false,
	},
	{
		Address: "elys1wd3463ejjkn32ae83e09vm65rgqkha45283cf4",
		Amount:  math.NewInt(1024900000),
		Claimed: false,
	},
	{
		Address: "elys1akv78wqjzvmyuyvaf93vss4fnh2utsfkheknn7",
		Amount:  math.NewInt(1014400000),
		Claimed: false,
	},
	{
		Address: "elys1mlszlgnvkmagtsdeaprr3cclgynzmjz38ldq2v",
		Amount:  math.NewInt(1013400000),
		Claimed: false,
	},
	{
		Address: "elys1y7zfelflwnm9aaldd9xp5x4p8azn8vddxzhem3",
		Amount:  math.NewInt(1011900000),
		Claimed: false,
	},
	{
		Address: "elys12k2an3fnmyxq7ll66q3aq3zmzkae29c6pv5f8n",
		Amount:  math.NewInt(1000000000),
		Claimed: false,
	},
	{
		Address: "elys1kd7n8nwgmedlnec30yqyv2ape29nwkt6vdulzy",
		Amount:  math.NewInt(1000000000),
		Claimed: false,
	},
	{
		Address: "elys1nz2qu2jdygcpnw7xjs3kxw4kr8dwug7h2nrp39",
		Amount:  math.NewInt(1000000000),
		Claimed: false,
	},
	{
		Address: "elys1w4pjkrel2quxx0wzmpmyx79nqwe8aqd59fuqv3",
		Amount:  math.NewInt(1000000000),
		Claimed: false,
	},
	{
		Address: "elys1aawzjnrelfytjgf3v06ws9hc9nq9hm0ssc53p7",
		Amount:  math.NewInt(998300000),
		Claimed: false,
	},
	{
		Address: "elys1zwlm9us6lpwvw4023x2y32w7799xdp8c95dagz",
		Amount:  math.NewInt(991800000),
		Claimed: false,
	},
	{
		Address: "elys1exqdj205n6ze74amfmw4evaq25nrtx59zh8mdg",
		Amount:  math.NewInt(990400000),
		Claimed: false,
	},
	{
		Address: "elys1a49thj2m764mh6pvwpthlmj70cdmhjeu2vm3ht",
		Amount:  math.NewInt(989600000),
		Claimed: false,
	},
	{
		Address: "elys1y8r34w6w0czzmzw3w30kadhxyqv0yn66c5mser",
		Amount:  math.NewInt(983400000),
		Claimed: false,
	},
	{
		Address: "elys1n2944shcnh9ee4ggwndmfc4kh3g4t3235j5y97",
		Amount:  math.NewInt(981700000),
		Claimed: false,
	},
	{
		Address: "elys1f7sr680szn8zrzxre9mt0c7umdalgh9wydt7hk",
		Amount:  math.NewInt(975700000),
		Claimed: false,
	},
	{
		Address: "elys1ltv32akcjtp7drqskvrajt03xr5pa079gakk3s",
		Amount:  math.NewInt(969600000),
		Claimed: false,
	},
	{
		Address: "elys124tjvtkxns5eqdwmaafh9pedemgzva6c40m2ra",
		Amount:  math.NewInt(961100000),
		Claimed: false,
	},
	{
		Address: "elys185e9eucf24myfe0fsyzxxxf2t2j88zacfatpk4",
		Amount:  math.NewInt(958800000),
		Claimed: false,
	},
	{
		Address: "elys1mu26eqsel6ft7zf4ce596nwcmp5umegjk6txpm",
		Amount:  math.NewInt(958100000),
		Claimed: false,
	},
	{
		Address: "elys1s5j7y6tnn7nlg884dl0agfv44a5tqjampyww2m",
		Amount:  math.NewInt(957100000),
		Claimed: false,
	},
	{
		Address: "elys1vvy8g6amqexhz9wn8vtepkgczp7a5c7alt8u5p",
		Amount:  math.NewInt(956100000),
		Claimed: false,
	},
	{
		Address: "elys1ezjdtssv7uwl77k392f7jdyqed2gc5gmas3l9y",
		Amount:  math.NewInt(955700000),
		Claimed: false,
	},
	{
		Address: "elys1xqvauc9myhlr87g80usn3vdsc5rdv5g0gdkdp4",
		Amount:  math.NewInt(953200000),
		Claimed: false,
	},
	{
		Address: "elys1h0z3ccv5y5g4kdnyqyu7ncv7f5w6f922hefqxa",
		Amount:  math.NewInt(947100000),
		Claimed: false,
	},
	{
		Address: "elys1dcy4edtqk8wpp7kn23tg50lvfvj3dpcq2zncg3",
		Amount:  math.NewInt(943700000),
		Claimed: false,
	},
	{
		Address: "elys17fapjthm3t27anz59808cy7y9uuypcnk5zv927",
		Amount:  math.NewInt(940700000),
		Claimed: false,
	},
	{
		Address: "elys1907xdvead672jf9n8m395yxcdhsnu5vfeswkuw",
		Amount:  math.NewInt(933000000),
		Claimed: false,
	},
	{
		Address: "elys106tw4h2fze0lxwjxyk3rj3jdl086dqzkmj0s59",
		Amount:  math.NewInt(918100000),
		Claimed: false,
	},
	{
		Address: "elys1ydrskl8w27n3dyxzgl233uazpgxarerdemjw7g",
		Amount:  math.NewInt(908800000),
		Claimed: false,
	},
	{
		Address: "elys17ctsn0rtdsygg3hsnz24wp0f7alph6d7a25dqm",
		Amount:  math.NewInt(903300000),
		Claimed: false,
	},
	{
		Address: "elys1cprnkgx2wmgzlxs6d3lknr30etyqmn665769wp",
		Amount:  math.NewInt(896500000),
		Claimed: false,
	},
	{
		Address: "elys1kg7djfj2guwzjlpvwvt5vqk6nv7maxg8304kw7",
		Amount:  math.NewInt(891300000),
		Claimed: false,
	},
	{
		Address: "elys15sx364xeelrcjqlcfmjs0dy6plfwmj4flstsc0",
		Amount:  math.NewInt(881900000),
		Claimed: false,
	},
	{
		Address: "elys12cw08yg7mm5l8fdg4aa6yf0tfy7ehmgasdh758",
		Amount:  math.NewInt(864700000),
		Claimed: false,
	},
	{
		Address: "elys1g43662df0afgd7q6xvn8eqkpcsyuw3krkz7azt",
		Amount:  math.NewInt(864300000),
		Claimed: false,
	},
	{
		Address: "elys1dxtpwk53ck7934xtn7ly38pf3fyjy5axtltvgq",
		Amount:  math.NewInt(862800000),
		Claimed: false,
	},
	{
		Address: "elys1rmwxnv2jt6wpvypyr9zf42q508tcj5hujmpstl",
		Amount:  math.NewInt(853100000),
		Claimed: false,
	},
	{
		Address: "elys18uy24ttg5fdruzq6732txjwxyh5lmq63xs5uag",
		Amount:  math.NewInt(852000000),
		Claimed: false,
	},
	{
		Address: "elys1h295kur87swge9sjurpeqnpwy86d87406huj7j",
		Amount:  math.NewInt(848300000),
		Claimed: false,
	},
	{
		Address: "elys1dhg6n69c5wzn95p8nflc8xu7rss0c7d6q54hvl",
		Amount:  math.NewInt(846200000),
		Claimed: false,
	},
	{
		Address: "elys1qdut9w6swxtn0cpnaud2qv6psh5y0suw2ys6kl",
		Amount:  math.NewInt(843300000),
		Claimed: false,
	},
	{
		Address: "elys17sx8n05ts8ylqzm6cj09mndtftqan824jztekr",
		Amount:  math.NewInt(841200000),
		Claimed: false,
	},
	{
		Address: "elys1mvfcavgjg4mjx2wamlzw0nzy25ahea7zspjv0u",
		Amount:  math.NewInt(839000000),
		Claimed: false,
	},
	{
		Address: "elys1upclff2d88xah68mgyquayfd7klp4udm7wf79c",
		Amount:  math.NewInt(830800000),
		Claimed: false,
	},
	{
		Address: "elys1sr36sm855wr0n4spspaknae7c64rdn09an7ft8",
		Amount:  math.NewInt(829800000),
		Claimed: false,
	},
	{
		Address: "elys1paslvj6znkyhpca2lugvhlsxkqtepjxxqmhajx",
		Amount:  math.NewInt(824000000),
		Claimed: false,
	},
	{
		Address: "elys10yqd3m8xxwdzd7s34gf5pat3h847sdelkcw7jt",
		Amount:  math.NewInt(821900000),
		Claimed: false,
	},
	{
		Address: "elys14va59m5pel7u7pav55ckm6zcer4ja4jtpqfya5",
		Amount:  math.NewInt(820100000),
		Claimed: false,
	},
	{
		Address: "elys1k0gngvg60gxhx0rezs3u07syuh7qajdq2un4pe",
		Amount:  math.NewInt(819700000),
		Claimed: false,
	},
	{
		Address: "elys1g8qynt25y7g87ah6dnwjmwt7t9fx6mjpxd00dy",
		Amount:  math.NewInt(813800000),
		Claimed: false,
	},
	{
		Address: "elys1hel9a046a74rrqemqrh4xfj3wyesspru5lxqnj",
		Amount:  math.NewInt(802800000),
		Claimed: false,
	},
	{
		Address: "elys1lt4styturu2q5qhgrpllq97cwhjh943k5s832l",
		Amount:  math.NewInt(798400000),
		Claimed: false,
	},
	{
		Address: "elys1zu65wnlv7s5rj0jdtqqz77390cfelr5dr9epq0",
		Amount:  math.NewInt(794800000),
		Claimed: false,
	},
	{
		Address: "elys1dq38sahlh9mtn3j4j4khl423w5m8dw8sldmx2v",
		Amount:  math.NewInt(786400000),
		Claimed: false,
	},
	{
		Address: "elys10w4mcmmpkz87rc40ctgl4cysjeyh63fcaqjgj9",
		Amount:  math.NewInt(783900000),
		Claimed: false,
	},
	{
		Address: "elys1fxu0vhtd608petrrx9pwjfv2c0zxxnh5q3w00w",
		Amount:  math.NewInt(782600000),
		Claimed: false,
	},
	{
		Address: "elys1hthcymqwvymcp0uhhy6s0psmqeq72e2ewppdqy",
		Amount:  math.NewInt(782300000),
		Claimed: false,
	},
	{
		Address: "elys19stygl5wsg3xc8hrkt6rgjuqec47t7cnm7m3tx",
		Amount:  math.NewInt(780100000),
		Claimed: false,
	},
	{
		Address: "elys1hrpddxpasr9pp5qu0w0w9us3g3v8hmvcrmwr4f",
		Amount:  math.NewInt(778500000),
		Claimed: false,
	},
	{
		Address: "elys14wpl8ylafdl90stlvtr6cpc867zxu5a94ur5sa",
		Amount:  math.NewInt(774100000),
		Claimed: false,
	},
	{
		Address: "elys1j9drjscm4ma5tuutcs0qhl4lhawrtyqx8vv9ge",
		Amount:  math.NewInt(769600000),
		Claimed: false,
	},
	{
		Address: "elys18wukalcm9rzkhufe22qdy7a640eumx8a3j3ey9",
		Amount:  math.NewInt(763800000),
		Claimed: false,
	},
	{
		Address: "elys1feur3exaph0k52qrs8cuzj6x73c92a3jhp5kkf",
		Amount:  math.NewInt(763100000),
		Claimed: false,
	},
	{
		Address: "elys1ve84qzcpra0l6hzdzr9k29j3yg6kjwup47hk7x",
		Amount:  math.NewInt(763000000),
		Claimed: false,
	},
	{
		Address: "elys1tfhdesnqxhlg66yw78q9y7gclkwsh0u37x7s78",
		Amount:  math.NewInt(759300000),
		Claimed: false,
	},
	{
		Address: "elys1dklqe0jzyh4vp94t5g3jmwsrgs3tgez5zavtgy",
		Amount:  math.NewInt(757200000),
		Claimed: false,
	},
	{
		Address: "elys1jr0rr0ldajr3ce7n9x73z7ff4njy77mfx3up6a",
		Amount:  math.NewInt(757100000),
		Claimed: false,
	},
	{
		Address: "elys1sgftl47xr8unupg30mcvvtx698dyv4eveqxhye",
		Amount:  math.NewInt(755100000),
		Claimed: false,
	},
	{
		Address: "elys1qm8ttsqm99uucnjqg7h2qzlrl6rxcay73xcv3s",
		Amount:  math.NewInt(754100000),
		Claimed: false,
	},
	{
		Address: "elys1xu0hnklz74pvge0rjw93xrvx06r5h0esz6fxfh",
		Amount:  math.NewInt(753800000),
		Claimed: false,
	},
	{
		Address: "elys1vxmq9ssrunnrfy0cq8vmm7wxypnucupdulqg49",
		Amount:  math.NewInt(753700000),
		Claimed: false,
	},
	{
		Address: "elys1dsnfjekjj23vtuemd8as5jmtf99ql5yh5nwn3f",
		Amount:  math.NewInt(751900000),
		Claimed: false,
	},
	{
		Address: "elys1f44qd2yw4f007dqddltptaqyqly9qs05a8ud8p",
		Amount:  math.NewInt(750200000),
		Claimed: false,
	},
	{
		Address: "elys1afusmz8753gre3z22xpvzfgxwx2tp8se8cmsyj",
		Amount:  math.NewInt(745700000),
		Claimed: false,
	},
	{
		Address: "elys16lt3fswfcux9txcdmdvpk0fdedsglhxx0gm9vw",
		Amount:  math.NewInt(745000000),
		Claimed: false,
	},
	{
		Address: "elys16fudeln2gt8yzjxyun4htk50xpppengv3x9u08",
		Amount:  math.NewInt(741000000),
		Claimed: false,
	},
	{
		Address: "elys17aaurv2k8nmtvr3mffq8xt9h26j0296ax35uvv",
		Amount:  math.NewInt(738600000),
		Claimed: false,
	},
	{
		Address: "elys1gz25lck5a7hk9ryk0t4gwthp0ctn7cl73a43w6",
		Amount:  math.NewInt(737000000),
		Claimed: false,
	},
	{
		Address: "elys1xcufk0vf0fhwddwcqfw66f67kjdpc4daxm3k5e",
		Amount:  math.NewInt(736000000),
		Claimed: false,
	},
	{
		Address: "elys18wh6r5escah3798kafny6vkr8s8pc3fcewwazh",
		Amount:  math.NewInt(734200000),
		Claimed: false,
	},
	{
		Address: "elys142ufxszy0t6xr43u7yza3fj5scadyag040decn",
		Amount:  math.NewInt(733500000),
		Claimed: false,
	},
	{
		Address: "elys17khs75h8rhnvtvlwrfa84tcm2sdmw8thaxnufx",
		Amount:  math.NewInt(732500000),
		Claimed: false,
	},
	{
		Address: "elys1rzm4xtz48rd9khv9yay2dkvjvfymvct6tas9gq",
		Amount:  math.NewInt(728500000),
		Claimed: false,
	},
	{
		Address: "elys1qczrusl58d9kselv6mp50dyajat86z48g6pvz7",
		Amount:  math.NewInt(728100000),
		Claimed: false,
	},
	{
		Address: "elys13rtl5gylg99vp9l2868sn6nny069c6cvkn32l5",
		Amount:  math.NewInt(723900000),
		Claimed: false,
	},
	{
		Address: "elys1q7qxzstrqp7gx7n898yn9gnxkuemk3kh9k0ajz",
		Amount:  math.NewInt(721500000),
		Claimed: false,
	},
	{
		Address: "elys1yen5f0ej9njg0d9pa8nn2hwjpqqm36zj843ms3",
		Amount:  math.NewInt(711200000),
		Claimed: false,
	},
	{
		Address: "elys198d6jvzfehc203kap2fd7jyvk5n2x8q9ye6lkw",
		Amount:  math.NewInt(708100000),
		Claimed: false,
	},
	{
		Address: "elys1adut4rgeurw5cgqhpa083sklzvf6s3pw4z3hc8",
		Amount:  math.NewInt(703200000),
		Claimed: false,
	},
	{
		Address: "elys19f5s5jjj262jf8h6jxyv7hj4wrq8pz0h7ayl6c",
		Amount:  math.NewInt(701400000),
		Claimed: false,
	},
	{
		Address: "elys1c7hwa7999v0em0rxhj403j8gng9h0xq7jjeh8v",
		Amount:  math.NewInt(701000000),
		Claimed: false,
	},
	{
		Address: "elys18cwccrn7q0ntrkzwxnh9eltt7rdhwgca96lgd6",
		Amount:  math.NewInt(700700000),
		Claimed: false,
	},
	{
		Address: "elys1nm2s7g4p5wf0q9y5vvvhj2cfwedvqea3mt0qh6",
		Amount:  math.NewInt(699500000),
		Claimed: false,
	},
	{
		Address: "elys1h2ay8awsa26zftfyjdc2f2y00uywfqflp28y6w",
		Amount:  math.NewInt(698800000),
		Claimed: false,
	},
	{
		Address: "elys1x5nfacmsujthl26v26qwexdcacfcxy855wmrls",
		Amount:  math.NewInt(697800000),
		Claimed: false,
	},
	{
		Address: "elys1yyv8gvfjg2x9d6zqrpr849mt87k49swehlmglj",
		Amount:  math.NewInt(695100000),
		Claimed: false,
	},
	{
		Address: "elys1gkhzup5rxt5nutkct0fffe92y93w5ukzeypmfk",
		Amount:  math.NewInt(693400000),
		Claimed: false,
	},
	{
		Address: "elys1q8mjtc9q9vdujzqhnz6s3e2cz02q5ge5cn5xkd",
		Amount:  math.NewInt(693100000),
		Claimed: false,
	},
	{
		Address: "elys1cw56n0hsffdku878x8x2c66u0690p9j0rhjyhx",
		Amount:  math.NewInt(690100000),
		Claimed: false,
	},
	{
		Address: "elys1927v6a30jk8xwj6ggvdlht5efcqvu6l96u07xx",
		Amount:  math.NewInt(688100000),
		Claimed: false,
	},
	{
		Address: "elys1saw7chmy8aqzq7cye0ruk6xaf3d7ep3lfltugk",
		Amount:  math.NewInt(683500000),
		Claimed: false,
	},
	{
		Address: "elys1v8cj8e0wt5vdsen0wk56888kehgvh0mdv897vl",
		Amount:  math.NewInt(681000000),
		Claimed: false,
	},
	{
		Address: "elys1hmgpengvmgc4w4wr03hrmpcqluhf78qmkr08cd",
		Amount:  math.NewInt(680500000),
		Claimed: false,
	},
	{
		Address: "elys1vqurd5tsn293d3j55dxynknrd73cc7p2fl9g5p",
		Amount:  math.NewInt(678000000),
		Claimed: false,
	},
	{
		Address: "elys1chn7mkkfx57hawxp7yfk642xjh69c8uqnfjvyq",
		Amount:  math.NewInt(677100000),
		Claimed: false,
	},
	{
		Address: "elys1rtljztddr3phdmdzr46mwdkyhlgcnkca485u2m",
		Amount:  math.NewInt(676700000),
		Claimed: false,
	},
	{
		Address: "elys1mvz5cq9d2vr2gjgm26p8jpqn8p45e5r7fk0aee",
		Amount:  math.NewInt(675100000),
		Claimed: false,
	},
	{
		Address: "elys1p5de4a0mjwt83j8vhyfvznhuj3334fnaf6y2mr",
		Amount:  math.NewInt(671500000),
		Claimed: false,
	},
	{
		Address: "elys1yrdgzjjdvnyv23cv7m27ww882yhemx47vpf8uf",
		Amount:  math.NewInt(671500000),
		Claimed: false,
	},
	{
		Address: "elys1fd6hrwastmer296tzwvk76eztf54q93ufr0ar7",
		Amount:  math.NewInt(670500000),
		Claimed: false,
	},
	{
		Address: "elys1z4nvyjapw4aux3y7ndkcd86prekamfut2a3fpw",
		Amount:  math.NewInt(670300000),
		Claimed: false,
	},
	{
		Address: "elys14kkq5m86zxd5z9thanzzku7z74rc0d3zmyaef2",
		Amount:  math.NewInt(668400000),
		Claimed: false,
	},
	{
		Address: "elys1ztr0wzkqpprphya6aake797nvx5etppv6m3hed",
		Amount:  math.NewInt(668100000),
		Claimed: false,
	},
	{
		Address: "elys1kwmk8mfrna308nm0jftrgravuuwuuk3fyqyvxh",
		Amount:  math.NewInt(667700000),
		Claimed: false,
	},
	{
		Address: "elys1gsyssnqfjqefwhcj895lql98uuumtt2wjhvqd2",
		Amount:  math.NewInt(664600000),
		Claimed: false,
	},
	{
		Address: "elys1rf9yuf834s2n5wyware0qgq4x7mggw93vvxrzj",
		Amount:  math.NewInt(664200000),
		Claimed: false,
	},
	{
		Address: "elys184y9g23mhd99c37kl0fvcsvx5hs8pfswv3x2ep",
		Amount:  math.NewInt(664000000),
		Claimed: false,
	},
	{
		Address: "elys1sk8jueqvjktd2ttnxnf0wza0hjpxe0uwu6a9ky",
		Amount:  math.NewInt(659600000),
		Claimed: false,
	},
	{
		Address: "elys1n8jhupgsfg2hds9d4tmprxsgw3raandn600zmn",
		Amount:  math.NewInt(659100000),
		Claimed: false,
	},
	{
		Address: "elys170cw4szqvk0l33n3f7yngr7p4zkzwkcnuzarxu",
		Amount:  math.NewInt(658500000),
		Claimed: false,
	},
	{
		Address: "elys1sat45e2e5tq6yf2jntulsuc4xz9dv7rruwdv3x",
		Amount:  math.NewInt(658400000),
		Claimed: false,
	},
	{
		Address: "elys1zmwecv9ah7lpq5g4k94mmttxuqqf6m97rhf9rp",
		Amount:  math.NewInt(658400000),
		Claimed: false,
	},
	{
		Address: "elys18g6ah3rd4uy52tuegftdly0jyedscm04nh7gf8",
		Amount:  math.NewInt(653500000),
		Claimed: false,
	},
	{
		Address: "elys16hrehrsupllep64g0afkm72ll7cqwseclmecm5",
		Amount:  math.NewInt(652300000),
		Claimed: false,
	},
	{
		Address: "elys1htq0cs4jaaszthjvegkjrv822cr4v95vt3asqa",
		Amount:  math.NewInt(650900000),
		Claimed: false,
	},
	{
		Address: "elys1567kw90t9pvp46k8c6vtx3lqdc6ysqz766k89w",
		Amount:  math.NewInt(650700000),
		Claimed: false,
	},
	{
		Address: "elys1p7x9ajy32jnkpdtezqcjy96s38ymkhkdph4vua",
		Amount:  math.NewInt(649100000),
		Claimed: false,
	},
	{
		Address: "elys1tg4lhwc0rqrja53nl5ksxuwvax2yefvhjxcmmg",
		Amount:  math.NewInt(648600000),
		Claimed: false,
	},
	{
		Address: "elys17mmpk8m6s57k60e357en52gens88eukmmzja5q",
		Amount:  math.NewInt(647100000),
		Claimed: false,
	},
	{
		Address: "elys1yaxlhvxtlwz4l08wavkpuxwqw2km59v40vnkeu",
		Amount:  math.NewInt(646900000),
		Claimed: false,
	},
	{
		Address: "elys1gd87qucyk69anz9cllmaylztcpk7jkwywld0fu",
		Amount:  math.NewInt(646400000),
		Claimed: false,
	},
	{
		Address: "elys1nldu4d4eelpqwn9uxf09tegeafze9rrznum2vf",
		Amount:  math.NewInt(645500000),
		Claimed: false,
	},
	{
		Address: "elys1ehhqvhgamke0avc44u88ycjfpu0rcxycf08wep",
		Amount:  math.NewInt(644000000),
		Claimed: false,
	},
	{
		Address: "elys1tvrrzrukz9s2tuzupnel0h3t497mzu9hh54fwc",
		Amount:  math.NewInt(643000000),
		Claimed: false,
	},
	{
		Address: "elys1gc2qsw4ds8jjygfxxl7k53ezhvfw666aa95xe7",
		Amount:  math.NewInt(642800000),
		Claimed: false,
	},
	{
		Address: "elys1sc4dh5npu2ddswkf7gn62sudjjzzt86qrdp46x",
		Amount:  math.NewInt(642400000),
		Claimed: false,
	},
	{
		Address: "elys1jcxf482t5hqyx0gqsfae7yz7p2jsekvyhhv534",
		Amount:  math.NewInt(640600000),
		Claimed: false,
	},
	{
		Address: "elys1rasja0l29rd5lgzlqj33nrwn605dp8lyk80ucu",
		Amount:  math.NewInt(640200000),
		Claimed: false,
	},
	{
		Address: "elys10kpq0paxzcqqdrjkvqjs44m60ea0jfd6glt3uc",
		Amount:  math.NewInt(638800000),
		Claimed: false,
	},
	{
		Address: "elys1fgmq6jd35wc4qcjvd3wevftw5m22j5t34wd3d7",
		Amount:  math.NewInt(638800000),
		Claimed: false,
	},
	{
		Address: "elys1pcceyzn84t74uqrtag4pja5qulrnxqvx7cyytt",
		Amount:  math.NewInt(638600000),
		Claimed: false,
	},
	{
		Address: "elys18tzh3xzynkuy7ehwqey34mmff873zye5hupgzh",
		Amount:  math.NewInt(636700000),
		Claimed: false,
	},
	{
		Address: "elys19cuk0pu4wc8wnfsx4zsszcsz2een066c0f666w",
		Amount:  math.NewInt(636500000),
		Claimed: false,
	},
	{
		Address: "elys16rqw5q4j80y7m5wppekym9rylwkkuwm995klqm",
		Amount:  math.NewInt(635900000),
		Claimed: false,
	},
	{
		Address: "elys1txtlppw3aepa5pfasrdpj4ww9zrc8ds9ugqc4f",
		Amount:  math.NewInt(635000000),
		Claimed: false,
	},
	{
		Address: "elys1ztr4rdt2uefgqgjnd4e2ztyy3966hknnefkvzl",
		Amount:  math.NewInt(634700000),
		Claimed: false,
	},
	{
		Address: "elys1tapvntnssnm4vqqqh8jxnv0fswc87nh9gncvrn",
		Amount:  math.NewInt(633700000),
		Claimed: false,
	},
	{
		Address: "elys19k2y3lrlxprlca0x69pur2m98zw6k5d8pgnr4c",
		Amount:  math.NewInt(633500000),
		Claimed: false,
	},
	{
		Address: "elys1ddqpvnymyv0g2axcef88u3wd6fpfx7lyev4rhq",
		Amount:  math.NewInt(631400000),
		Claimed: false,
	},
	{
		Address: "elys1y6q6ws40raqs2zrnqxa9p4x3hxq5rn3g599pl4",
		Amount:  math.NewInt(630100000),
		Claimed: false,
	},
	{
		Address: "elys16v23a56268xf2whl6pf35epe9qh5vynttn8z7e",
		Amount:  math.NewInt(629300000),
		Claimed: false,
	},
	{
		Address: "elys15s3usen8dvw4kyvl9fuje82mspf6svkm2hn3aa",
		Amount:  math.NewInt(629200000),
		Claimed: false,
	},
	{
		Address: "elys1ja5fcujdygck7xlgr0r9qjqq4gl8mw5hhzky68",
		Amount:  math.NewInt(628500000),
		Claimed: false,
	},
	{
		Address: "elys10ulfxmkysgw7dly2v5zky0nfusp8shnj63memw",
		Amount:  math.NewInt(628000000),
		Claimed: false,
	},
	{
		Address: "elys1n23xedafmunpmr65ahe08754x03tgdz2e45hpd",
		Amount:  math.NewInt(626500000),
		Claimed: false,
	},
	{
		Address: "elys1e7hgmdgma8sntnepale0w6efdqnxmannz7yuep",
		Amount:  math.NewInt(625400000),
		Claimed: false,
	},
	{
		Address: "elys18zzs354zcvmp4kgg9ut0lza3wrtgtajqrpd3nn",
		Amount:  math.NewInt(624900000),
		Claimed: false,
	},
	{
		Address: "elys1whejvdne3cj6p5pnap8cnvgnlxdyd4dxwum69s",
		Amount:  math.NewInt(624800000),
		Claimed: false,
	},
	{
		Address: "elys1jt87q06ct8hfxza60vrxevrhtepvjpzwc4tfwx",
		Amount:  math.NewInt(624300000),
		Claimed: false,
	},
	{
		Address: "elys160uyt76cewx8g38mnktqpp9rj96d3clry009va",
		Amount:  math.NewInt(622400000),
		Claimed: false,
	},
	{
		Address: "elys1j06f2hxezhpafsdfmt37ghnf6lvz7p0qw66guz",
		Amount:  math.NewInt(622400000),
		Claimed: false,
	},
	{
		Address: "elys1kyetkvtl8d74dsefdljnhwc77axf7283vz3whv",
		Amount:  math.NewInt(616900000),
		Claimed: false,
	},
	{
		Address: "elys1kmk7z7hkc6ndpcedteck4548fnn3nexw4f5347",
		Amount:  math.NewInt(616500000),
		Claimed: false,
	},
	{
		Address: "elys1hfcsxncn9gdexhttqxt3d47cagx2enx9n9pfvf",
		Amount:  math.NewInt(615700000),
		Claimed: false,
	},
	{
		Address: "elys1jc7ffzkz2mp6p23nf0lsrm7ggwq9ynpm79cu8r",
		Amount:  math.NewInt(615600000),
		Claimed: false,
	},
	{
		Address: "elys1rxxqlr7qj2fmwlz5e7v4ymxmpzspzst2lnh8nz",
		Amount:  math.NewInt(615400000),
		Claimed: false,
	},
	{
		Address: "elys1t3pwtk5g34q8qu52hrmatake75w57hslkus05u",
		Amount:  math.NewInt(615400000),
		Claimed: false,
	},
	{
		Address: "elys1pp79dv2u2m5sjf0esh8nd08ujquvvv9trg2jyx",
		Amount:  math.NewInt(614300000),
		Claimed: false,
	},
	{
		Address: "elys15gmyy5ekx7z6purzfftkya8x9sug4wg6hd9fj5",
		Amount:  math.NewInt(613200000),
		Claimed: false,
	},
	{
		Address: "elys1w7de3ejt83mgj3keqkvq9p96s0qjygpxqkr6qd",
		Amount:  math.NewInt(612200000),
		Claimed: false,
	},
	{
		Address: "elys10cnrppgzgzw5z57uwpa2e8v536r5v29xa5r7z9",
		Amount:  math.NewInt(611100000),
		Claimed: false,
	},
	{
		Address: "elys1a6c8jlv9jmlmvs68my33n73flzhz8den5dzdx5",
		Amount:  math.NewInt(611100000),
		Claimed: false,
	},
	{
		Address: "elys1928ccsfn5ehs6capwdfmu9eyep56fav43thhxv",
		Amount:  math.NewInt(610200000),
		Claimed: false,
	},
	{
		Address: "elys1y97cppg563cc3d070dvsgf5cdmx6uu9ta6njw8",
		Amount:  math.NewInt(610000000),
		Claimed: false,
	},
	{
		Address: "elys1448zdtyt4p3atr5xrm3hqrj9s0tre3dqydkeja",
		Amount:  math.NewInt(609300000),
		Claimed: false,
	},
	{
		Address: "elys1ehvtqgqs5snpwa93mmw0rkjptcr09fghstluk0",
		Amount:  math.NewInt(609000000),
		Claimed: false,
	},
	{
		Address: "elys1ajjpa969qlekh25qzclx2gjzjupdt7cqenje2n",
		Amount:  math.NewInt(607500000),
		Claimed: false,
	},
	{
		Address: "elys19ltl6597kz7ke6kjm424ny3tcpvuqmrt2r2zj2",
		Amount:  math.NewInt(607100000),
		Claimed: false,
	},
	{
		Address: "elys18esdx9sf8n44dw97n3htyrjw6jnqmkzuz7xz2q",
		Amount:  math.NewInt(605800000),
		Claimed: false,
	},
	{
		Address: "elys1xe5x8u4hlxq2mewymwdet6e8j7g0jf86q26767",
		Amount:  math.NewInt(605800000),
		Claimed: false,
	},
	{
		Address: "elys1x7mnqj3f3z58evztngmkhztje728nxlt78ra57",
		Amount:  math.NewInt(603300000),
		Claimed: false,
	},
	{
		Address: "elys1qppd6jc6mtnhfa7cge2qqhgxpd0ylrml7neh3n",
		Amount:  math.NewInt(602600000),
		Claimed: false,
	},
	{
		Address: "elys16rvrnqsdckr8zdvd5gngjvy75kve24z5u0pjcn",
		Amount:  math.NewInt(601700000),
		Claimed: false,
	},
	{
		Address: "elys17dy5htsncg7mjyf68ecqqangrrgqjst48c7v4u",
		Amount:  math.NewInt(599900000),
		Claimed: false,
	},
	{
		Address: "elys1dhu8grdkd6sw9rw848kkfuf7e0x7a47eh0lzx2",
		Amount:  math.NewInt(599500000),
		Claimed: false,
	},
	{
		Address: "elys1w0hwp29v8cpcmqskwt9d6fxle0zsu2vgujtr5r",
		Amount:  math.NewInt(599500000),
		Claimed: false,
	},
	{
		Address: "elys1can64mxtqa3wzknsuzsspjlcvd0yd4panrtfyl",
		Amount:  math.NewInt(599400000),
		Claimed: false,
	},
	{
		Address: "elys16af03s965vqs4flkdyr4s45eyp7r4v6dmuju0e",
		Amount:  math.NewInt(599300000),
		Claimed: false,
	},
	{
		Address: "elys12pt55nz633ylfd5h2ze62unum4ddwupypwr9dv",
		Amount:  math.NewInt(598600000),
		Claimed: false,
	},
	{
		Address: "elys13v3sehsjmma5hj6maylj993ydhls39cew0trc5",
		Amount:  math.NewInt(598500000),
		Claimed: false,
	},
	{
		Address: "elys1da3gtqa3ljam2hs5pppe88ahgyg2kqknpng38y",
		Amount:  math.NewInt(598300000),
		Claimed: false,
	},
	{
		Address: "elys1lgjqa37h6ulyufcngzhuxcw6fuhuyguqehpzu5",
		Amount:  math.NewInt(598300000),
		Claimed: false,
	},
	{
		Address: "elys12l7drlyudsww0sep85gfcj3n3xf6r208n6re4l",
		Amount:  math.NewInt(597000000),
		Claimed: false,
	},
	{
		Address: "elys1wwze387w898fyuwarturspzjue8pmtn7yf0522",
		Amount:  math.NewInt(596500000),
		Claimed: false,
	},
	{
		Address: "elys1k2j4e73gqfzuza7l9hng48qht04gkju88ugk9l",
		Amount:  math.NewInt(595100000),
		Claimed: false,
	},
	{
		Address: "elys1e0v83jk03v04n2we9ha4hq4dwt7sq6gkg9xs5z",
		Amount:  math.NewInt(594500000),
		Claimed: false,
	},
	{
		Address: "elys175j5h8v4n9s6sekmdvmhrj4893cjplh9t9nghu",
		Amount:  math.NewInt(594200000),
		Claimed: false,
	},
	{
		Address: "elys1e9xkater59sg4g3tz0a9yvqery0erpsqrk7pfm",
		Amount:  math.NewInt(593900000),
		Claimed: false,
	},
	{
		Address: "elys1uhn7axp7fv5ttuh822dq08xk9r5vv6g8jghvwq",
		Amount:  math.NewInt(593900000),
		Claimed: false,
	},
	{
		Address: "elys1p4xaea2gs95pte28e7lzccf7l9jsxz2wm7nngx",
		Amount:  math.NewInt(592400000),
		Claimed: false,
	},
	{
		Address: "elys19z663dtlpeumx8uh8gz0d3ffjsh5wjydap50za",
		Amount:  math.NewInt(592300000),
		Claimed: false,
	},
	{
		Address: "elys1uhz32nd9trwqxe25v3xqux9pvcpw7nncs3wa0c",
		Amount:  math.NewInt(592000000),
		Claimed: false,
	},
	{
		Address: "elys1sj3aq3x33sadq2rvk455wfv6wdsqwam5k552gz",
		Amount:  math.NewInt(591400000),
		Claimed: false,
	},
	{
		Address: "elys14xtscjdd7vue6fz7rk75n0sfjtga5fg2ssmzn6",
		Amount:  math.NewInt(591100000),
		Claimed: false,
	},
	{
		Address: "elys1es4fz09x9eweqxjvymdtcgtnxxuej0ut4q64vp",
		Amount:  math.NewInt(590800000),
		Claimed: false,
	},
	{
		Address: "elys19d878g4ynzfxst6px5xk2fd2tqgu89dx86l0cl",
		Amount:  math.NewInt(590500000),
		Claimed: false,
	},
	{
		Address: "elys18fdu02ufsmak66z4qnm6u6es3ztkhdzep9q88v",
		Amount:  math.NewInt(589500000),
		Claimed: false,
	},
	{
		Address: "elys199k45lxgz5u558uk59jq0axwf87audgec22xkg",
		Amount:  math.NewInt(586800000),
		Claimed: false,
	},
	{
		Address: "elys1dlkr3zq5fvcdda95jpduuvd0eql6943t2xm3md",
		Amount:  math.NewInt(583900000),
		Claimed: false,
	},
	{
		Address: "elys1jr47lp2smyrj05f0w978pdw3thrqlsnv28h4qq",
		Amount:  math.NewInt(583800000),
		Claimed: false,
	},
	{
		Address: "elys1625t69vgndj86k3zn4cqyru3lz5uagmewr8fjk",
		Amount:  math.NewInt(581500000),
		Claimed: false,
	},
	{
		Address: "elys1uwcepjfa5q9fqu8fzhjr7l8zd96m2vpg4ue8ew",
		Amount:  math.NewInt(581100000),
		Claimed: false,
	},
	{
		Address: "elys1ua5wud6y8cqgat84zu64g3q7kkafq4efq3e8k5",
		Amount:  math.NewInt(580500000),
		Claimed: false,
	},
	{
		Address: "elys1qhh94kahtf5rj5m536xc7cs9tkheud0wta2r32",
		Amount:  math.NewInt(579600000),
		Claimed: false,
	},
	{
		Address: "elys1w894e9j6x4lng9snsmup6xut4fr07ky35ufnnl",
		Amount:  math.NewInt(579400000),
		Claimed: false,
	},
	{
		Address: "elys1fnk6nj24aa320mty5lf2l4mqvjjx82lujnepp2",
		Amount:  math.NewInt(576600000),
		Claimed: false,
	},
	{
		Address: "elys1va2xq5d9lj5v8aakxawgne3chn9phzwlh4dq99",
		Amount:  math.NewInt(576600000),
		Claimed: false,
	},
	{
		Address: "elys1xr3d3hdzc654jlfawlqynqk4htpyxlewm9ywg5",
		Amount:  math.NewInt(574800000),
		Claimed: false,
	},
	{
		Address: "elys1cy6ydatymvpza3j37daj836xs6hc5lxkcgh7xt",
		Amount:  math.NewInt(574600000),
		Claimed: false,
	},
	{
		Address: "elys1x2wcr92430nj4cnxx874dtqhacq84pdzlgj4rj",
		Amount:  math.NewInt(572900000),
		Claimed: false,
	},
	{
		Address: "elys1jlw0fszfhh2q85scydy55dsrw3ferjjyyetdl6",
		Amount:  math.NewInt(571600000),
		Claimed: false,
	},
	{
		Address: "elys1scqvstk4wlpunyz5c808rl0n4v24jr0f9qc6eg",
		Amount:  math.NewInt(569600000),
		Claimed: false,
	},
	{
		Address: "elys1xr93j93czlfeahg6gwx2l9ht6560dn5vd3f57h",
		Amount:  math.NewInt(566500000),
		Claimed: false,
	},
	{
		Address: "elys1u4j90gp5wu9q3m89hne8aznst7drrxng3p6ejz",
		Amount:  math.NewInt(565800000),
		Claimed: false,
	},
	{
		Address: "elys1j3es6w62xwj0ad8lg09eu3hjayv9cljycdsrry",
		Amount:  math.NewInt(565700000),
		Claimed: false,
	},
	{
		Address: "elys1t64nvz2mx8h9yj7dpe3l9tpwm5er4lkugurm0k",
		Amount:  math.NewInt(564900000),
		Claimed: false,
	},
	{
		Address: "elys1hwhtst9252jnzlfvqlrfxwept7dv0uwzllk2zz",
		Amount:  math.NewInt(564800000),
		Claimed: false,
	},
	{
		Address: "elys122vtl62w6k0dnenppy2x27segds0x9t7u3ssrz",
		Amount:  math.NewInt(560100000),
		Claimed: false,
	},
	{
		Address: "elys1n7suk7lzuam5yhlvdvlsmejn3t3r6n763t43ed",
		Amount:  math.NewInt(559900000),
		Claimed: false,
	},
	{
		Address: "elys17k22my0ut4qumcpslteuv32956r4jmypwa9rug",
		Amount:  math.NewInt(559300000),
		Claimed: false,
	},
	{
		Address: "elys1tey7xelcu07z6ctwpxzg3gm3d7jcn2quax3k7s",
		Amount:  math.NewInt(558900000),
		Claimed: false,
	},
	{
		Address: "elys18dsr6z7w8xw2wrgephceq6v0udnezj0rph62ed",
		Amount:  math.NewInt(556600000),
		Claimed: false,
	},
	{
		Address: "elys1a0gu6ul05w67f3lmem8lmlym9xjk5hjt0lq7vc",
		Amount:  math.NewInt(556400000),
		Claimed: false,
	},
	{
		Address: "elys1aud6egphphz94velufrqlpyyjqwe4ue2zhdx6u",
		Amount:  math.NewInt(554900000),
		Claimed: false,
	},
	{
		Address: "elys170vahwshlvd0ulw47m34uzvwklpm8t7ct3l2u5",
		Amount:  math.NewInt(554000000),
		Claimed: false,
	},
	{
		Address: "elys1xukpeq3q0s8gehu9xucs2fr7fzlqtqdludecqe",
		Amount:  math.NewInt(553800000),
		Claimed: false,
	},
	{
		Address: "elys1nv4vn7hq4ln32jpe0qw0vsv8mmkp9durzsq7dn",
		Amount:  math.NewInt(553600000),
		Claimed: false,
	},
	{
		Address: "elys1uxg9lhvetcsxayes9yq6jjqaykympfz8k6u33n",
		Amount:  math.NewInt(550200000),
		Claimed: false,
	},
	{
		Address: "elys154qggy89597vszvypep78ws0t6n284weqlazt7",
		Amount:  math.NewInt(549600000),
		Claimed: false,
	},
	{
		Address: "elys1695hjrhknrkmlkmxud8tuelzlxw4zj0hwmex7x",
		Amount:  math.NewInt(549600000),
		Claimed: false,
	},
	{
		Address: "elys1yvc5rwrdvqwpkzwyexucqnekcx5v93z4rpkzz6",
		Amount:  math.NewInt(549200000),
		Claimed: false,
	},
	{
		Address: "elys1j575t9askwsk66hjsqp96lndkct37cv0wywpcg",
		Amount:  math.NewInt(548800000),
		Claimed: false,
	},
	{
		Address: "elys1amac29dufxmkc9mzxl84xtfkzru9gyj7zr0qjv",
		Amount:  math.NewInt(546700000),
		Claimed: false,
	},
	{
		Address: "elys12j48ak4u0rrst690zcxj3gfqudvnjudsn33wxl",
		Amount:  math.NewInt(544000000),
		Claimed: false,
	},
	{
		Address: "elys1vwqnwkq3te9e0cdrck50sgafx7z3ayrlc4mzdv",
		Amount:  math.NewInt(543000000),
		Claimed: false,
	},
	{
		Address: "elys1rh40h6ljpwzt2f5r67tlph76hq2qcr622jczky",
		Amount:  math.NewInt(538000000),
		Claimed: false,
	},
	{
		Address: "elys1lkhsxela3kn5da6a52uv9fc7yxu5awr8hjfl0q",
		Amount:  math.NewInt(536700000),
		Claimed: false,
	},
	{
		Address: "elys1vxu3mear5jmw0sq5r9qshge4pke453nnuqv7j3",
		Amount:  math.NewInt(534700000),
		Claimed: false,
	},
	{
		Address: "elys1prr2x6gy692cllg2zxtuh5pz70n2a7tj88wq8l",
		Amount:  math.NewInt(534299999),
		Claimed: false,
	},
	{
		Address: "elys147s0w5xsh6mwr8dut9gtjvq0f8skvemnac4srv",
		Amount:  math.NewInt(533200000),
		Claimed: false,
	},
	{
		Address: "elys1j7cm0qhs4tvjpu4fyxd4sqzrhwc892a3lcxe7a",
		Amount:  math.NewInt(533200000),
		Claimed: false,
	},
	{
		Address: "elys1js638yd8j0yjwg3md9qm2aqc0zc008h0c3m3af",
		Amount:  math.NewInt(531600000),
		Claimed: false,
	},
	{
		Address: "elys168z3uhn9kwpps4j7erwh2nl5z20ugxqf3rkshd",
		Amount:  math.NewInt(531200000),
		Claimed: false,
	},
	{
		Address: "elys1yt8ksrraqwuve8q2ep8nuhl09n9ga99h4kgqzs",
		Amount:  math.NewInt(529600000),
		Claimed: false,
	},
	{
		Address: "elys1rfrfnkjc3t9zh9hn96p327cjyytk48fcha4kt0",
		Amount:  math.NewInt(528400000),
		Claimed: false,
	},
	{
		Address: "elys1vaywwa4aehwh3m4z5qex3qnl0c4ek7f8u43jqe",
		Amount:  math.NewInt(527900000),
		Claimed: false,
	},
	{
		Address: "elys1850emzx7y2pe2d74kuk0sh56c6jtu8nv6zprxn",
		Amount:  math.NewInt(527600000),
		Claimed: false,
	},
	{
		Address: "elys1t6t55q2l04zq6x37czruadhyaxsn5k02h38xu7",
		Amount:  math.NewInt(527600000),
		Claimed: false,
	},
	{
		Address: "elys1rn53wgd3jmghpqv3dvuf8s2d6ahmme8udg6cx3",
		Amount:  math.NewInt(526299999),
		Claimed: false,
	},
	{
		Address: "elys1ux4ej0zefpz54jhe7lu3ya868yrlsgcyzp4f7k",
		Amount:  math.NewInt(525799999),
		Claimed: false,
	},
	{
		Address: "elys1a6az3eg9paeltlpp98etmyl3ny9vgmlxxnagrk",
		Amount:  math.NewInt(525299999),
		Claimed: false,
	},
	{
		Address: "elys1hclrtf2e8p9nehg0ww456afw4tcza8pw2ef2jz",
		Amount:  math.NewInt(524000000),
		Claimed: false,
	},
	{
		Address: "elys1n705mtwx02wn7rz3f3h6j05v7y67j44kvzaws9",
		Amount:  math.NewInt(521900000),
		Claimed: false,
	},
	{
		Address: "elys1s3g756xdpxfnk8cg63ar5e0e9kmacfl4w8vrkg",
		Amount:  math.NewInt(517299999),
		Claimed: false,
	},
	{
		Address: "elys1flgktq77pprdqtrudypwm0qzl4ne9dsyptxxqs",
		Amount:  math.NewInt(512900000),
		Claimed: false,
	},
	{
		Address: "elys168xfluxe0psdzgvl0r77thz6h4sms3lk34ts9z",
		Amount:  math.NewInt(511400000),
		Claimed: false,
	},
	{
		Address: "elys15m5nuhuyxd0wenfhjcdvjgjgdy6jmtm0fpje5e",
		Amount:  math.NewInt(508000000),
		Claimed: false,
	},
	{
		Address: "elys1ldtlv60kn4qecmdedekpfpqtm46pnnc039jph3",
		Amount:  math.NewInt(504400000),
		Claimed: false,
	},
	{
		Address: "elys1xyrx3d267t2qr2h4xhl9j2nhpamtene2g4myg3",
		Amount:  math.NewInt(503800000),
		Claimed: false,
	},
	{
		Address: "elys1tmxynhy3hxrxtwauxl7rzxde7svpscy5y0vsjs",
		Amount:  math.NewInt(501900000),
		Claimed: false,
	},
	{
		Address: "elys1cmnnpqdj5j3sr560nzzl8avrl5c0z3ky2ptzjp",
		Amount:  math.NewInt(501500000),
		Claimed: false,
	},
	{
		Address: "elys16pm0wankqenwgfuva3jarxlthaylesxz3gpt87",
		Amount:  math.NewInt(497400000),
		Claimed: false,
	},
	{
		Address: "elys1mv3emx60h9je6g2td37zl0cp3gwc3qn09glp43",
		Amount:  math.NewInt(496800000),
		Claimed: false,
	},
	{
		Address: "elys124drgqmp60zngeq6mt7wxr6exhgl53ed4ma26z",
		Amount:  math.NewInt(496600000),
		Claimed: false,
	},
	{
		Address: "elys180t046h26kxkwlkmpq4kk43lnt6grur0x4vuuk",
		Amount:  math.NewInt(496500000),
		Claimed: false,
	},
	{
		Address: "elys1clw0rr9z36lc9z8rvkyzt0uqwjx83d4sl49mjq",
		Amount:  math.NewInt(496500000),
		Claimed: false,
	},
	{
		Address: "elys1x2xuzqlsnxpmg6p79sdpp6955a2z468lmtvxys",
		Amount:  math.NewInt(496300000),
		Claimed: false,
	},
	{
		Address: "elys1wrk5cfp9ajc0rx0mpkqxn6hm6z3em0wg3r79ur",
		Amount:  math.NewInt(495900000),
		Claimed: false,
	},
	{
		Address: "elys1y6ks8y8u4fjqlk6j2g2e3tlp6nh794wgxvr7x0",
		Amount:  math.NewInt(495800000),
		Claimed: false,
	},
	{
		Address: "elys16tyavrstn6jnuweu6dsqlmuk0y9fqx26t203l9",
		Amount:  math.NewInt(495200000),
		Claimed: false,
	},
	{
		Address: "elys1g37l5pcrprn2krdv04kqxtsg8x8r9vqx3p996p",
		Amount:  math.NewInt(495000000),
		Claimed: false,
	},
	{
		Address: "elys1qk74379nc34ynvc4gmw2zrwlscn62mx5x66raz",
		Amount:  math.NewInt(490400000),
		Claimed: false,
	},
	{
		Address: "elys17rxmjsxl7rf2ftnyqgag993qy5eqltmgtz5ydy",
		Amount:  math.NewInt(488600000),
		Claimed: false,
	},
	{
		Address: "elys1lkdhkt0vq5cwtt4khe90ur4cnjetutz3qdnf3f",
		Amount:  math.NewInt(484600000),
		Claimed: false,
	},
	{
		Address: "elys1fwmcetqmn2re4dfcq9gh3xzhsj0vc2mqmyayl2",
		Amount:  math.NewInt(483700000),
		Claimed: false,
	},
	{
		Address: "elys1h4pycycmjh7ye4u028l2gu3pqcrjhgq8k8syht",
		Amount:  math.NewInt(483200000),
		Claimed: false,
	},
	{
		Address: "elys1twg50yuldq0zyt5e3eevlq2jv6jln2e96kmt7l",
		Amount:  math.NewInt(481000000),
		Claimed: false,
	},
	{
		Address: "elys1249qt88r98ujrqwc0h7fz36udepkswst03vpvl",
		Amount:  math.NewInt(479800000),
		Claimed: false,
	},
	{
		Address: "elys1c52wvffdqmzdzglp6ky9v70e68xnjt4yak88ax",
		Amount:  math.NewInt(479700000),
		Claimed: false,
	},
	{
		Address: "elys1q0qqhq6vlxlp4efx259cd3kmfawtqkfteg9stm",
		Amount:  math.NewInt(477900000),
		Claimed: false,
	},
	{
		Address: "elys1sg7yt4l763qjx9x0vj9f4gqm3r6qv5fw5xsdew",
		Amount:  math.NewInt(476100000),
		Claimed: false,
	},
	{
		Address: "elys1t3ajryk6y5ljfk85yxhkut73f3ex0e42nd0af7",
		Amount:  math.NewInt(476100000),
		Claimed: false,
	},
	{
		Address: "elys1hsgyqxg3azzc7nwx4zv6upcgmhmx9fxcrjx3us",
		Amount:  math.NewInt(475200000),
		Claimed: false,
	},
	{
		Address: "elys1phymv9fz63yxcn3wjp94w04sxjazamhce9j0av",
		Amount:  math.NewInt(475000000),
		Claimed: false,
	},
	{
		Address: "elys1xsnlq9mzrj40p2cya8pvupjqgquaatg7mxqse5",
		Amount:  math.NewInt(474800000),
		Claimed: false,
	},
	{
		Address: "elys17645ef6rynx2x5uztgjfjq5m62k3q0cux7vtzd",
		Amount:  math.NewInt(473200000),
		Claimed: false,
	},
	{
		Address: "elys1778rw4t8fhcjm7xrrue4j0q44uqw0yrpkdlwwp",
		Amount:  math.NewInt(473100000),
		Claimed: false,
	},
	{
		Address: "elys1lv2w9yzyme4xvm3jf4tjhrvtp37plvf9l4s403",
		Amount:  math.NewInt(473100000),
		Claimed: false,
	},
	{
		Address: "elys1svlfdw0tn73h4klqzpkf60z20z3np3ddctmyyw",
		Amount:  math.NewInt(473100000),
		Claimed: false,
	},
	{
		Address: "elys122z8076n2cglyytaupve0c7jtrkugqdfjldw8p",
		Amount:  math.NewInt(471700000),
		Claimed: false,
	},
	{
		Address: "elys1n2csdcq6252cjtjcq997msravh5r7vh2a2m5nh",
		Amount:  math.NewInt(470400000),
		Claimed: false,
	},
	{
		Address: "elys13hq3x8m54vfv7dt8g87flhycjajg5297qd28zt",
		Amount:  math.NewInt(469900000),
		Claimed: false,
	},
	{
		Address: "elys1knjw5300mjmmfx7h6ktvqf7nnyuygk6lfw3jcx",
		Amount:  math.NewInt(469900000),
		Claimed: false,
	},
	{
		Address: "elys1rfu4ua3avdwm8ate3hczw7x3r2a86fzjjfsp43",
		Amount:  math.NewInt(469400000),
		Claimed: false,
	},
	{
		Address: "elys1xsn2hu68kfayxn8283e5gf3vrkp7qyn43zenwk",
		Amount:  math.NewInt(464100000),
		Claimed: false,
	},
	{
		Address: "elys1wmyhjssyyd05q2nglydg0fptwycz00d2dgg5hm",
		Amount:  math.NewInt(461600000),
		Claimed: false,
	},
	{
		Address: "elys187a2y0d0xzkpa33l8dz3u3wlhkyfyflm8nggxr",
		Amount:  math.NewInt(459900000),
		Claimed: false,
	},
	{
		Address: "elys1wuc72nzmdfn67q7mz4jajyu3xpmcj9mk3qznpl",
		Amount:  math.NewInt(457700000),
		Claimed: false,
	},
	{
		Address: "elys1kpkktn24dhuwntz88r3pteuk47n7uf8qtnk9kh",
		Amount:  math.NewInt(454800000),
		Claimed: false,
	},
	{
		Address: "elys1u9rx32gsk8wvv6x54j5fy899gfm5v5ckw25upx",
		Amount:  math.NewInt(454700000),
		Claimed: false,
	},
	{
		Address: "elys1t3uum0npc9caf862yq8zetxj5d2rzl42hjt6dq",
		Amount:  math.NewInt(454600000),
		Claimed: false,
	},
	{
		Address: "elys1pamf6cxspvf3e6lyu53e030xgl7v5ljmynj9r0",
		Amount:  math.NewInt(453600000),
		Claimed: false,
	},
	{
		Address: "elys1qkm4khavrzpyrlrcmvn3jvzfl43676kvcq6new",
		Amount:  math.NewInt(451700000),
		Claimed: false,
	},
	{
		Address: "elys17c97ykhyfefykezpw38j9tah45hjw4ulr34tdx",
		Amount:  math.NewInt(449900000),
		Claimed: false,
	},
	{
		Address: "elys1jsl39dn5swpa6rmx3lfpydyg72dg0dqgzatjsk",
		Amount:  math.NewInt(449600000),
		Claimed: false,
	},
	{
		Address: "elys1ctwp6y9herj4terh69tyj7h0mflnlqsvevgnus",
		Amount:  math.NewInt(449200000),
		Claimed: false,
	},
	{
		Address: "elys14cqqk6kv2z55nsht5rgpaec9spq04t9f8qz0ek",
		Amount:  math.NewInt(447600000),
		Claimed: false,
	},
	{
		Address: "elys1dyvadsu0txtjevpupt9g030wfmmd7np6f3frxu",
		Amount:  math.NewInt(447600000),
		Claimed: false,
	},
	{
		Address: "elys1ygqcg2v3qej563xtcwxtkjjnna096gt7uc2f6v",
		Amount:  math.NewInt(446900000),
		Claimed: false,
	},
	{
		Address: "elys1xt30vlqpfa6sutqk4lml5e7kpydfx4he8fy32z",
		Amount:  math.NewInt(446500000),
		Claimed: false,
	},
	{
		Address: "elys1pj542zfqux96www8g0t9wnsa7j93m4w5fjmm44",
		Amount:  math.NewInt(446300000),
		Claimed: false,
	},
	{
		Address: "elys1vl3e5wcarytmtmuwry5c8nrla4syd8ujwtwudq",
		Amount:  math.NewInt(439500000),
		Claimed: false,
	},
	{
		Address: "elys10h5qkujhw8jnztppngtjtg5wy2ma9xwqte5c6z",
		Amount:  math.NewInt(439000000),
		Claimed: false,
	},
	{
		Address: "elys18c6e3ap85emlddz32f75u0kdf5jrk9ced5k9m9",
		Amount:  math.NewInt(432600000),
		Claimed: false,
	},
	{
		Address: "elys1v2kn3v2fq8msa650sg9ftv329x4kcm62q3pr6j",
		Amount:  math.NewInt(429900000),
		Claimed: false,
	},
	{
		Address: "elys1ajgymx5xtwc4u0pfdndvjqvsl7t5rkqf0crcch",
		Amount:  math.NewInt(427900000),
		Claimed: false,
	},
	{
		Address: "elys1j3uknvx6rm3cu2mpqgxsad6jhe65u7pffzcavw",
		Amount:  math.NewInt(427300000),
		Claimed: false,
	},
	{
		Address: "elys15sz727cdzsgepw2aqc5slq0qj43ladjupw0au6",
		Amount:  math.NewInt(425500000),
		Claimed: false,
	},
	{
		Address: "elys1xn3f3vm7cnkfz5c9tqzrr2e9krzdcw7jkgtv7f",
		Amount:  math.NewInt(425200000),
		Claimed: false,
	},
	{
		Address: "elys1ferh26prrn8ncc0malj7zrytkgxlmljmrxsj9z",
		Amount:  math.NewInt(424800000),
		Claimed: false,
	},
	{
		Address: "elys1c88xrxjanlth68qfvqqwpkpxvfl89kqnzgayad",
		Amount:  math.NewInt(423700000),
		Claimed: false,
	},
	{
		Address: "elys18fltum8llunux5lyqrcr9l37def82cpzmr96p9",
		Amount:  math.NewInt(423000000),
		Claimed: false,
	},
	{
		Address: "elys1l69el4xqjrkrjsalg9zhtl3spym5870f5mjpnk",
		Amount:  math.NewInt(421400000),
		Claimed: false,
	},
	{
		Address: "elys1wl2dcwm5r89xlwqx0grzd5wgkja8nktp9ffn8t",
		Amount:  math.NewInt(421300000),
		Claimed: false,
	},
	{
		Address: "elys10hvz04hh92xzct5hxnpsn5h2fp3p4ammuh90rx",
		Amount:  math.NewInt(418300000),
		Claimed: false,
	},
	{
		Address: "elys10az8y677yhvjgehq52y5nkhy3u5vsua9rw8mcu",
		Amount:  math.NewInt(418200000),
		Claimed: false,
	},
	{
		Address: "elys1ckwppdpjgmux2e9knsrj6m703f46tvn0f9cpf0",
		Amount:  math.NewInt(417900000),
		Claimed: false,
	},
	{
		Address: "elys1sderc2e5v54ku6ys6rygnsnm3lldtnxw73a46t",
		Amount:  math.NewInt(417700000),
		Claimed: false,
	},
	{
		Address: "elys1d5leksg2ueyt886w6g3z20yqffn34scsv66574",
		Amount:  math.NewInt(416300000),
		Claimed: false,
	},
	{
		Address: "elys1nwmz9acleuagpzu40mpcjcsm8xye87t7w0y8yq",
		Amount:  math.NewInt(415300000),
		Claimed: false,
	},
	{
		Address: "elys1dd6ndzgpv7vuneg0fvw0psvg865ssdxfu2wm7v",
		Amount:  math.NewInt(414000000),
		Claimed: false,
	},
	{
		Address: "elys15g3ls5t2ar6gst7plszut0yy6k9urfjta3ucwq",
		Amount:  math.NewInt(412900000),
		Claimed: false,
	},
	{
		Address: "elys1r4r7m6xyvaa0td0nklm4y85xf9cc04kk6elhn9",
		Amount:  math.NewInt(412700000),
		Claimed: false,
	},
	{
		Address: "elys1n0z409ndz4kwla4q2mdxm8nd282lx5dnt6yq07",
		Amount:  math.NewInt(410400000),
		Claimed: false,
	},
	{
		Address: "elys17v8nctuam9gn0sp4a9jdh4sy6eshdt9lgw7yux",
		Amount:  math.NewInt(409300000),
		Claimed: false,
	},
	{
		Address: "elys180cj9v695r4qkk3e6k7gcj5zurx674dyr8wd0k",
		Amount:  math.NewInt(408500000),
		Claimed: false,
	},
	{
		Address: "elys1dd27v99564tfr2ctahjapakxwa4djng428mz5p",
		Amount:  math.NewInt(408500000),
		Claimed: false,
	},
	{
		Address: "elys17e3ftc8mu826add5paj3xxgc425zu48796dk49",
		Amount:  math.NewInt(408100000),
		Claimed: false,
	},
	{
		Address: "elys1gsyjls3rmnx0k7mzvl5shwmdqcmtu968nvn2w4",
		Amount:  math.NewInt(405700000),
		Claimed: false,
	},
	{
		Address: "elys1lm25s78k6jpfq6kgg4lmd7hfukcyyrvxy9na60",
		Amount:  math.NewInt(405300000),
		Claimed: false,
	},
	{
		Address: "elys1x8ygmpj5dl6x682d085flr2u8avevzpqdfjgt3",
		Amount:  math.NewInt(405100000),
		Claimed: false,
	},
	{
		Address: "elys1vzdk6ehzjz4pu4ztlu7fkl6hzqjccvsydll8t5",
		Amount:  math.NewInt(404100000),
		Claimed: false,
	},
	{
		Address: "elys1hs96kexuc66whw9cx9kve2957j8n6ag02eqvzy",
		Amount:  math.NewInt(403900000),
		Claimed: false,
	},
	{
		Address: "elys1w0qnyd4jef6tlcj7nusd8py2es5vk6k8p0dcgq",
		Amount:  math.NewInt(402700000),
		Claimed: false,
	},
	{
		Address: "elys17pmymeg4hxcpy68gauqyzs96t7xnx8xu6wylss",
		Amount:  math.NewInt(401400000),
		Claimed: false,
	},
	{
		Address: "elys1qfksff087n26zfcnlncfumzuzrrcxcjqeejdkl",
		Amount:  math.NewInt(399900000),
		Claimed: false,
	},
	{
		Address: "elys197azlmrp8wvz4446jfnxu605snegznl4akqkdz",
		Amount:  math.NewInt(399500000),
		Claimed: false,
	},
	{
		Address: "elys14pdkkwpj7gvtjvcr27lqp908w0cjgwettlv8uh",
		Amount:  math.NewInt(399400000),
		Claimed: false,
	},
	{
		Address: "elys1ed58f6wavepkxhl4h4jyuljt5u8jmhrwvffyfg",
		Amount:  math.NewInt(399200000),
		Claimed: false,
	},
	{
		Address: "elys1fmsn9apdpgh0dg32v746rf7xdckf2fgnv77gmq",
		Amount:  math.NewInt(397400000),
		Claimed: false,
	},
	{
		Address: "elys1x693aqslyv3znhltlfuqckwjq6nml5urvethny",
		Amount:  math.NewInt(397300000),
		Claimed: false,
	},
	{
		Address: "elys1xx64nu0xgxeefstl5v72nlgsy8f73w0tmnvmza",
		Amount:  math.NewInt(397300000),
		Claimed: false,
	},
	{
		Address: "elys1zux0djv2xasrpddung50estdghmcddtayp45cf",
		Amount:  math.NewInt(397300000),
		Claimed: false,
	},
	{
		Address: "elys1xn70hs6xspp5yn68upkddwkc2kmfhpw8kncy20",
		Amount:  math.NewInt(396900000),
		Claimed: false,
	},
	{
		Address: "elys1dpy72pqepeagpgs9x4uv8dj4040a0t8ckqwftg",
		Amount:  math.NewInt(396200000),
		Claimed: false,
	},
	{
		Address: "elys17tysp87mgxrugc4gfdu388ut34krmra5s7pk54",
		Amount:  math.NewInt(395700000),
		Claimed: false,
	},
	{
		Address: "elys1pat0nhud2zygpkev98zykmfksf93vga5zs6u5s",
		Amount:  math.NewInt(395200000),
		Claimed: false,
	},
	{
		Address: "elys14ekp40hfhjxj9lvx3zpj5qc7xavprczf3jx2y0",
		Amount:  math.NewInt(394700000),
		Claimed: false,
	},
	{
		Address: "elys19mrfdeq4gl6fx5srsu2w8a8f9vcxha8e29cnc7",
		Amount:  math.NewInt(394600000),
		Claimed: false,
	},
	{
		Address: "elys17ftetnw5fw8mgywa39dwczuzwhp3cwvpavhsy2",
		Amount:  math.NewInt(390600000),
		Claimed: false,
	},
	{
		Address: "elys10xeay3ghpul80hkqeua020qhrljjg74ky0kux0",
		Amount:  math.NewInt(389900000),
		Claimed: false,
	},
	{
		Address: "elys1vkerl4d9uzuahlkxlfjd4dzhn97hkf9nn7r748",
		Amount:  math.NewInt(389800000),
		Claimed: false,
	},
	{
		Address: "elys14p6mj6a9entjsfsugmthrj7ucggru8dcge6jjh",
		Amount:  math.NewInt(387500000),
		Claimed: false,
	},
	{
		Address: "elys10tc7vx4jqdakzg9n8vd0gqqp4n6ven43eg9t7n",
		Amount:  math.NewInt(385400000),
		Claimed: false,
	},
	{
		Address: "elys1ch32xm2qwln3q8p7qc6rwvcnschwfzs9zqtx67",
		Amount:  math.NewInt(385400000),
		Claimed: false,
	},
	{
		Address: "elys1dcnntha63n2sqq0gl8ar6kgm5lqwzqsmn2km65",
		Amount:  math.NewInt(385100000),
		Claimed: false,
	},
	{
		Address: "elys12uegyk6vjk0rfpqp7uhmwmmvg8de6g7ttyy4n0",
		Amount:  math.NewInt(384900000),
		Claimed: false,
	},
	{
		Address: "elys1zqv58alxz2cyvdqdq2nz393kjy56edu3wsfffm",
		Amount:  math.NewInt(384500000),
		Claimed: false,
	},
	{
		Address: "elys1wryt0g4g7h8eldzankzpx55kraljjepxyeej6t",
		Amount:  math.NewInt(383500000),
		Claimed: false,
	},
	{
		Address: "elys1jgwhexdj6j9r56e72mcfqg6sx8dljngc0ey0zk",
		Amount:  math.NewInt(382800000),
		Claimed: false,
	},
	{
		Address: "elys1ctq3fwh3wrfrfzajzjj4n3pljx90gu902relud",
		Amount:  math.NewInt(382300000),
		Claimed: false,
	},
	{
		Address: "elys1x67n2q0fp2g2uf09xtprxlved879sr6al5cwa0",
		Amount:  math.NewInt(380700000),
		Claimed: false,
	},
	{
		Address: "elys1dqst5kxwruc4g4fflv9rejmrcm0ezkxcvxdjfk",
		Amount:  math.NewInt(378500000),
		Claimed: false,
	},
	{
		Address: "elys1p7dsm5hnwh5gxwj6kr0hspnu7wkl9t5rdafrrj",
		Amount:  math.NewInt(376400000),
		Claimed: false,
	},
	{
		Address: "elys1vudzp43y5cg4q9t7h79cc94nc952gdnh5tu0md",
		Amount:  math.NewInt(376200000),
		Claimed: false,
	},
	{
		Address: "elys10khyx0rvpta07me9eczqkt508ts90wg43v6gce",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys10mx9js8zkwhncg37202t6w7jr6hxj69lxutt40",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys12l5uufvm82xyt725t9kvm6d0qkpm4jmw9cq0ye",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys16ww0xyrh3cwr8s7pyt4mn797gxhvlvd3sm48r5",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys179urwhjmhuphzc4f58z4du988tlkysz0t3yued",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys17dlfyt3ek0chwapx24f60d5tz3n5s8qykxr9me",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys17lwha3f5jue86hwcfdd6crpxz85j8nf8m7hg5g",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys17x2a2j039uwv0jeetm3tgcs9yx09s6tvjz5rns",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys18ghky0yypgvpu64jcgt7ff02qkrs260urn95zr",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys18nsqnn7nnf4lwehx8e2se060fk6usf6xvdqen3",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys19rlvshhja3xkkc3n8ahzxxmffshn6cxwlw8s2m",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1angyrk8ydszqad2qzdqa3uqadfndtuvfphcmxw",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1cfkfel08mzdtzg4j88mfhf3cdsnztqy9v0deht",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1cmg7he3rtxemmdccq5tny0avrrzq4798dls6ge",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1csez43ha37qnmvsjea2d96are3eme4289jv3zj",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1dxy35pljr6npukye5ve3rnga3dhnra987usu7y",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1eqpmwmwacn38x6e6fsc7w9pzelmufavmea0emj",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1eyq5s8ke5c3m73j5wkw76jrx9gqrueppfau5zq",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1ghqyejcdr4sqjl8kg840n0hd3hlarkycn8580m",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1h7n08cjpdpvaj3hgfz7dlm8343tevtrytj6v2c",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1jq3terhmyuct8p6xatskg720fy66hkat8re8a6",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1kgmvfpaq9gnu8hnzwf9nsddtnnnvy6fwa8867l",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1kwkw6078xqqjkdurt85a33yrhmlanhqjetsvel",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1lzy5e8szs7lflgetq3rrkavq8npwcngrsqwwgn",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1me5fwxft0fg3cduvjmxq5uqmm46r6hfz8lfysu",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1ntmuee2yf8xpps0ljg0es0ef2ufefhh98pdjuh",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1p4wxarept0vhet0zyh4em7yjclql8c7jnlpnkq",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1qxcg9wtku4f0cl3wm5amfc3zcve27a2avw4gnr",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1r4w0aj05479h73m5j4zrcd4rzhmmwtl03jecvy",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1rgf3eheqxvr92mq5h2kpg5slralmw3rnywr0cc",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1spd6sfrkf4m6mx0rrdxqdyturpq6eqdfm3h08x",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1txczzt8eruvx6ut4r8h0vc8zdedyuzfaxc4ese",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1v0rdzrh8hty4ldxj0ec42pxyss5aad6ktqkuem",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1v8h045fa6vh00dmkk8wwsu3hktpw8p8nj580l3",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1vmxy52qjkx75kfdzcs4faufsxkr9fhs0rw5suu",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1ycqgj60a34wu6w3x3rfxj0xn6ulk06wlu0pudk",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1yy9m72m3wte2g4q42p2w2nup83gvrwe60munuc",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1tppd4j6l6fwtm5nyhqhnkpzuf94zv9rpc32djf",
		Amount:  math.NewInt(375000000),
		Claimed: false,
	},
	{
		Address: "elys1v2hsm69lhk33qt8qrm2pumqcwmfs0nr56dleal",
		Amount:  math.NewInt(372900000),
		Claimed: false,
	},
	{
		Address: "elys1eugurlxpg3rxwty6ymq7dulr4ffyuqyt870tj4",
		Amount:  math.NewInt(372100000),
		Claimed: false,
	},
	{
		Address: "elys1px6qrmraz2k4jjskaufu9kefyvqs8htlc0xfts",
		Amount:  math.NewInt(371100000),
		Claimed: false,
	},
	{
		Address: "elys1cd46uzct3ej8r7vw6yhy86gx5fytu3x0n02l5t",
		Amount:  math.NewInt(370000000),
		Claimed: false,
	},
	{
		Address: "elys1kqh8fjrpj2cqmunfx7fmz0rxfyx06l3f829e3k",
		Amount:  math.NewInt(370000000),
		Claimed: false,
	},
	{
		Address: "elys10fwa5zs9pml7n7st65n3fz8ef25ggmwaqqk6x9",
		Amount:  math.NewInt(369000000),
		Claimed: false,
	},
	{
		Address: "elys16c26mgp6k7r5l4h8p9n9mnx486rwxftjacyqnd",
		Amount:  math.NewInt(368900000),
		Claimed: false,
	},
	{
		Address: "elys1vngrm9sfajyckv7jc57k6p0kdxxzcvtghgge6u",
		Amount:  math.NewInt(368900000),
		Claimed: false,
	},
	{
		Address: "elys1hxghx7ms8g2zgc4j087y60xl7u8vgmzss73yzy",
		Amount:  math.NewInt(368700000),
		Claimed: false,
	},
	{
		Address: "elys19j77h6800j5m3qfq03xpkeh0m698z46u40f3wt",
		Amount:  math.NewInt(368200000),
		Claimed: false,
	},
	{
		Address: "elys12rchqus5pl9cud9zn6np7lhayerjrtcm2l7xh0",
		Amount:  math.NewInt(366500000),
		Claimed: false,
	},
	{
		Address: "elys1y9pttrvtt3ngdsxpkxzd7vk04g8em7klrmfjzc",
		Amount:  math.NewInt(364700000),
		Claimed: false,
	},
	{
		Address: "elys1kw6mw70wafdxgp2n8s4lscx04du8ka6dakhd9t",
		Amount:  math.NewInt(364000000),
		Claimed: false,
	},
	{
		Address: "elys1vqwyz34d3sm5uljsg52em7juxgwf029rk45p9n",
		Amount:  math.NewInt(361100000),
		Claimed: false,
	},
	{
		Address: "elys16yf30hy79u5mtl38652rjan0n5h2attc5634j0",
		Amount:  math.NewInt(359700000),
		Claimed: false,
	},
	{
		Address: "elys1adqeglychswah04f8w3jfnurjkdt9f8heueucw",
		Amount:  math.NewInt(359100000),
		Claimed: false,
	},
	{
		Address: "elys1t43hjxjnrcpkdsc46s72kh5av72mlsu696rtfv",
		Amount:  math.NewInt(357000000),
		Claimed: false,
	},
	{
		Address: "elys1zgs5nnwhpxg4sjyw7xerm4humluxua90hjcy9t",
		Amount:  math.NewInt(355500000),
		Claimed: false,
	},
	{
		Address: "elys1dx34096tpg0da8ng4ja9wh8s3rmqvah2dglkju",
		Amount:  math.NewInt(355100000),
		Claimed: false,
	},
	{
		Address: "elys1a9zmhn5ycxtz6h0x3f50cy2tdcve335dmd93d6",
		Amount:  math.NewInt(354600000),
		Claimed: false,
	},
	{
		Address: "elys1nsqef88pu5n2qa3gxppsm6ud0mua93xrqg3cm0",
		Amount:  math.NewInt(354600000),
		Claimed: false,
	},
	{
		Address: "elys1huapl52kveecspkmc0rgkk9jyef68fpdt59uzp",
		Amount:  math.NewInt(353600000),
		Claimed: false,
	},
	{
		Address: "elys1xfztc2framxxan8dnpcmza7rk47xkp7l7xfsx8",
		Amount:  math.NewInt(348200000),
		Claimed: false,
	},
	{
		Address: "elys1xx33y2exg8s4vjusr4v2gye3upv4ethdaadp2s",
		Amount:  math.NewInt(345500000),
		Claimed: false,
	},
	{
		Address: "elys15wnjx4z4knp0h9cfru3cp0txpy5n9g5p67pgjm",
		Amount:  math.NewInt(344900000),
		Claimed: false,
	},
	{
		Address: "elys1vs883vrt8xq92ltmrj7xqtjngjuj83l9qye8zr",
		Amount:  math.NewInt(344700000),
		Claimed: false,
	},
	{
		Address: "elys1jn8lcl4f0dzug93700a6nacmz07k492essja90",
		Amount:  math.NewInt(334400000),
		Claimed: false,
	},
	{
		Address: "elys1c2w7sxnp423935c47zwv32rgmk2zmgt8nx6egk",
		Amount:  math.NewInt(334100000),
		Claimed: false,
	},
	{
		Address: "elys12mlgdxag6fp9vrc664s9f6sdty3m88wa0p5xgh",
		Amount:  math.NewInt(330800000),
		Claimed: false,
	},
	{
		Address: "elys1azf9rzckh45zht5aql3lcpweqyq293qw4srcu6",
		Amount:  math.NewInt(329900000),
		Claimed: false,
	},
	{
		Address: "elys1dh04mrlvzzh6c5vaj90mugx42scjteq7l033r8",
		Amount:  math.NewInt(327500000),
		Claimed: false,
	},
	{
		Address: "elys1k3aqeumdf5ctt5qw4zjca40vdhm47h009d6skh",
		Amount:  math.NewInt(327200000),
		Claimed: false,
	},
	{
		Address: "elys179sxg8mvxfmpmuumfgc7jss4n0zh3mapr885vv",
		Amount:  math.NewInt(327000000),
		Claimed: false,
	},
	{
		Address: "elys1w6dpa2l2fc3v2rc9dpuhqg4lt6ykvlqgq00q8r",
		Amount:  math.NewInt(326800000),
		Claimed: false,
	},
	{
		Address: "elys1clhv0yaq7fqxp5dgrvd72r578jrs532v3vx574",
		Amount:  math.NewInt(325500000),
		Claimed: false,
	},
	{
		Address: "elys1ffcwyzfdzf2ym0sw05d8j9enxcn9dutjlnwmck",
		Amount:  math.NewInt(322200000),
		Claimed: false,
	},
	{
		Address: "elys1jrz7w5xww8t4ppt66pacrmdl0chwvm5e9t0qd4",
		Amount:  math.NewInt(321200000),
		Claimed: false,
	},
	{
		Address: "elys1sces8mevflwf74erdnle48eulc426e6lp2st4h",
		Amount:  math.NewInt(320200000),
		Claimed: false,
	},
	{
		Address: "elys1fs9gsnxl8y5d5fjt2za4m3nzgx73ecjvfnwaau",
		Amount:  math.NewInt(319800000),
		Claimed: false,
	},
	{
		Address: "elys1yw5kj5zjzqztuahflm75f3rav7cc7s0lr48fwk",
		Amount:  math.NewInt(318400000),
		Claimed: false,
	},
	{
		Address: "elys1syjj9uv3wsnhx7zrjh4pkunwnv75yrcvct8nus",
		Amount:  math.NewInt(315800000),
		Claimed: false,
	},
	{
		Address: "elys1jxu3vc5prjfqg2nxez9m6qnwg20wfcwa3vyl5l",
		Amount:  math.NewInt(315400000),
		Claimed: false,
	},
	{
		Address: "elys170df6q5tznl2922029vsnu7v9j6twszj2yug5x",
		Amount:  math.NewInt(315100000),
		Claimed: false,
	},
	{
		Address: "elys1423lrkcyufus67w7j79hyjfjr0hd35kxgc82s3",
		Amount:  math.NewInt(314400000),
		Claimed: false,
	},
	{
		Address: "elys1gjrznum06xkwzxans6dk6gs2v50zwxkq8rhwpc",
		Amount:  math.NewInt(313900000),
		Claimed: false,
	},
	{
		Address: "elys1xt6qyjq5puta53srr8hvgh39vf095u5ycfh994",
		Amount:  math.NewInt(313700000),
		Claimed: false,
	},
	{
		Address: "elys158g286vs8zh5t324dphes5877gg5nm4ktjlzfx",
		Amount:  math.NewInt(312400000),
		Claimed: false,
	},
	{
		Address: "elys1trcsja8udlrv28283vlujlstvnenumtafjgww6",
		Amount:  math.NewInt(311100000),
		Claimed: false,
	},
	{
		Address: "elys1ur2u2sq5e6dv8xpd3p58xaxxtd7qnh477r6n7y",
		Amount:  math.NewInt(310400000),
		Claimed: false,
	},
	{
		Address: "elys10wzftxrgecgzlm8er9agdevlxzy34utrrhx7p9",
		Amount:  math.NewInt(310300000),
		Claimed: false,
	},
	{
		Address: "elys1yh4cqydd5c5xnfy3xaenahrvzxeu5xqll0r8l3",
		Amount:  math.NewInt(310300000),
		Claimed: false,
	},
	{
		Address: "elys1gfjrv9u5dnrg25yv7qqmu9f88nxnzd6eam0v97",
		Amount:  math.NewInt(307200000),
		Claimed: false,
	},
	{
		Address: "elys1tkpcnd66w9cteh46h7kyf2azqrwwent0y4f6kr",
		Amount:  math.NewInt(306700000),
		Claimed: false,
	},
	{
		Address: "elys1w6cf5r94q0k8znh0sgpnzhkqlkk38fhzcp6sk0",
		Amount:  math.NewInt(306400000),
		Claimed: false,
	},
	{
		Address: "elys10zhjwdmn0c6v0kn80t5l0v2dnxak6m5j7vhman",
		Amount:  math.NewInt(305400000),
		Claimed: false,
	},
	{
		Address: "elys1c4u0h8h4wcwf3vq05ey3fee04tq445cz4xt2z4",
		Amount:  math.NewInt(305000000),
		Claimed: false,
	},
	{
		Address: "elys1xvp68gtr03v6yd8j4lrr8sj5p54y49exrz829k",
		Amount:  math.NewInt(303500000),
		Claimed: false,
	},
	{
		Address: "elys1mpah7dwgsukmkuyatgncltmfy8cgthm4mr2g8m",
		Amount:  math.NewInt(303000000),
		Claimed: false,
	},
	{
		Address: "elys1n652vc2y2tku4p6lfcydhtqwnrrz9pzlx43tk3",
		Amount:  math.NewInt(302100000),
		Claimed: false,
	},
	{
		Address: "elys1etcf9np7q9sc6lskgwyyue5m25qlt4e5na2dkd",
		Amount:  math.NewInt(301200000),
		Claimed: false,
	},
	{
		Address: "elys1062pqrt6z6w0nsp6t25ekxmgxx9dl03jmvajx2",
		Amount:  math.NewInt(301100000),
		Claimed: false,
	},
	{
		Address: "elys1uzucdnm59vl87522tdv40wl9lkskd7h2sp2a6k",
		Amount:  math.NewInt(298400000),
		Claimed: false,
	},
	{
		Address: "elys1t90wpnw279yfr0yx8p7ul0p79agampzsgh6lvt",
		Amount:  math.NewInt(298100000),
		Claimed: false,
	},
	{
		Address: "elys10xfytvupvnyufztwzztn8w3d0hasvddja2z2z8",
		Amount:  math.NewInt(297600000),
		Claimed: false,
	},
	{
		Address: "elys19866n97fe93up7xea2p6u5r25e6hgdj57a24j7",
		Amount:  math.NewInt(296700000),
		Claimed: false,
	},
	{
		Address: "elys1hl934chhhag7v6nunx0k98kz86u5y5xkaer555",
		Amount:  math.NewInt(296300000),
		Claimed: false,
	},
	{
		Address: "elys1tketjrth9xwnqxa0vj0mx7mdp000yek69dkhke",
		Amount:  math.NewInt(296200000),
		Claimed: false,
	},
	{
		Address: "elys1a7enjntdrytnxscm3469q4sad5a3uxug4z7f0p",
		Amount:  math.NewInt(295900000),
		Claimed: false,
	},
	{
		Address: "elys1k7pp0fmh0ev0dxtnt3d7kg035hpayxf39dtq44",
		Amount:  math.NewInt(295500000),
		Claimed: false,
	},
	{
		Address: "elys1u8kpj6l0h4j8mvl5lw6aag4xcceulqywkjucrw",
		Amount:  math.NewInt(295300000),
		Claimed: false,
	},
	{
		Address: "elys1sc0t39l8fpzmdw6l0lzzrwrw3d8486n3jgaruq",
		Amount:  math.NewInt(294600000),
		Claimed: false,
	},
	{
		Address: "elys1z2uyrzjkuld7vcf420m0rlyqugtj3vv68uj8fw",
		Amount:  math.NewInt(294400000),
		Claimed: false,
	},
	{
		Address: "elys15fw3np2yx9rtrexjcda7uwk78hh0lltcpzrna5",
		Amount:  math.NewInt(293100000),
		Claimed: false,
	},
	{
		Address: "elys1rtqje443v5vpnjl4xnnu42u80jaxymh3ukwsxj",
		Amount:  math.NewInt(292700000),
		Claimed: false,
	},
	{
		Address: "elys1egcuw75mhxv9sv2vudvnd9p4dywjw5aukwsqyh",
		Amount:  math.NewInt(291600000),
		Claimed: false,
	},
	{
		Address: "elys1lq6pyh0shpk8u5vgxnwcmfzkhjrf96kedh33pa",
		Amount:  math.NewInt(289800000),
		Claimed: false,
	},
	{
		Address: "elys1vx04r7mgd6ch3esvy05xzqnqc0rfvqwkyldyyz",
		Amount:  math.NewInt(287900000),
		Claimed: false,
	},
	{
		Address: "elys1ehyrladq2yq07500kkkuwt02w30g9qrkqtx2t2",
		Amount:  math.NewInt(286300000),
		Claimed: false,
	},
	{
		Address: "elys18fnfpvcj3pfr8cj2hfz4y0qlj4679f5zwy9rqz",
		Amount:  math.NewInt(285700000),
		Claimed: false,
	},
	{
		Address: "elys19mr0pm536jfr2ta8588eay49z94negwvakwxrp",
		Amount:  math.NewInt(284500000),
		Claimed: false,
	},
	{
		Address: "elys1acyqnvk2l6zxpd8le02wgtl783z0mvve63zyw8",
		Amount:  math.NewInt(282900000),
		Claimed: false,
	},
	{
		Address: "elys1u3nk4reaed7ntk3v8kkal42yl6aygcsxwymqd8",
		Amount:  math.NewInt(281300000),
		Claimed: false,
	},
	{
		Address: "elys1spe42m04gq93j7hj793lq7avlmehydm0nct58q",
		Amount:  math.NewInt(280700000),
		Claimed: false,
	},
	{
		Address: "elys1k0zdx92k6d60vszdhm45hk4tpn06f6kglr5tu9",
		Amount:  math.NewInt(280000000),
		Claimed: false,
	},
	{
		Address: "elys1vf68rkgk502k4ja6l2mq9t8pcdrp5qpzml53rg",
		Amount:  math.NewInt(278200000),
		Claimed: false,
	},
	{
		Address: "elys1wccs45ag53w0csr3wmc5cp5kedxzl464sujqx6",
		Amount:  math.NewInt(277900000),
		Claimed: false,
	},
	{
		Address: "elys18ql99h7usawt42lljha9x6yk8xhxghr3asqjjp",
		Amount:  math.NewInt(276900000),
		Claimed: false,
	},
	{
		Address: "elys17rep78t7d55e0ffdf3r4s0tws7rv4ec098f63h",
		Amount:  math.NewInt(274500000),
		Claimed: false,
	},
	{
		Address: "elys1vvuhmfk3qeu7lmsf0mel347slzylelu0vmhjya",
		Amount:  math.NewInt(271800000),
		Claimed: false,
	},
	{
		Address: "elys1ccfn84s8pnjzdv8tsfv22xd5j2fx2yekdle604",
		Amount:  math.NewInt(270700000),
		Claimed: false,
	},
	{
		Address: "elys19ykpuyausrhkhmet4gmjnl9qw5zhuc73k2kffj",
		Amount:  math.NewInt(270300000),
		Claimed: false,
	},
	{
		Address: "elys1cn7fqru45dejr696jxutys997nlxk690ekt283",
		Amount:  math.NewInt(268900000),
		Claimed: false,
	},
	{
		Address: "elys1p5rvpd5wzk4k0ngyp5u6t05p796s0tsk8dsh68",
		Amount:  math.NewInt(266899999),
		Claimed: false,
	},
	{
		Address: "elys1n84w0hcr6yc6l2pjss05deur7wa0lkdf2vjdua",
		Amount:  math.NewInt(264899999),
		Claimed: false,
	},
	{
		Address: "elys1s7nxxj6pnrakz0n29afcf9rtl8r3n2fz9k80qq",
		Amount:  math.NewInt(264600000),
		Claimed: false,
	},
	{
		Address: "elys1gvt4z647pks9mzd7k5zc07p03dmvsdne8ndndd",
		Amount:  math.NewInt(262600000),
		Claimed: false,
	},
	{
		Address: "elys1dt34kla8te3t3k0p53qvvmp0fgjaeujnthq7x3",
		Amount:  math.NewInt(262300000),
		Claimed: false,
	},
	{
		Address: "elys18meutp7u2jnksfjn2czwv9u47ed0h4zqk3kxn8",
		Amount:  math.NewInt(262000000),
		Claimed: false,
	},
	{
		Address: "elys1zph7lv5tj0grhmlym4fhv2u2cjrv6eg7xmugpz",
		Amount:  math.NewInt(261700000),
		Claimed: false,
	},
	{
		Address: "elys15u002gfste7xwmycfe9sjxsvxaadmelcfyu2wh",
		Amount:  math.NewInt(261399999),
		Claimed: false,
	},
	{
		Address: "elys1drt64enfrd27qkzzppzd3v824ka7rvq5qwtr89",
		Amount:  math.NewInt(261000000),
		Claimed: false,
	},
	{
		Address: "elys174s83ddz3lmzdaymj78jyxljtfjwuqy5g4vt7y",
		Amount:  math.NewInt(260800000),
		Claimed: false,
	},
	{
		Address: "elys1dnz744efyeymyc77e7qjwshg0s7my56x2nz6mx",
		Amount:  math.NewInt(260800000),
		Claimed: false,
	},
	{
		Address: "elys1xc76d2z9athg4amu92cg3xt8v95dgmantf59ew",
		Amount:  math.NewInt(260700000),
		Claimed: false,
	},
	{
		Address: "elys1j3dpast3jf8kw8mmgmz48j3nedv7yk6r9h5ftx",
		Amount:  math.NewInt(260600000),
		Claimed: false,
	},
	{
		Address: "elys17czz0ftq5w52hzcgjmwzx82u0a4d5fgn6sc84d",
		Amount:  math.NewInt(260300000),
		Claimed: false,
	},
	{
		Address: "elys129cla5847k5wg7x947qgdn09rk2aav70d7lzjg",
		Amount:  math.NewInt(259899999),
		Claimed: false,
	},
	{
		Address: "elys1vza3kkueeh0ka3f0zpzav8zqdnlk56tfgrp9rh",
		Amount:  math.NewInt(259600000),
		Claimed: false,
	},
	{
		Address: "elys1mscdau2xps35j6k9h059jg69mpnv33ks6ewrng",
		Amount:  math.NewInt(257800000),
		Claimed: false,
	},
	{
		Address: "elys1tv0z9z0uw24xqkklq8fjnvp825q2v0t29t6vt9",
		Amount:  math.NewInt(257700000),
		Claimed: false,
	},
	{
		Address: "elys1rq75r89nlu6cujfv0y7axfqtdh4nmlz8lvqugq",
		Amount:  math.NewInt(257000000),
		Claimed: false,
	},
	{
		Address: "elys14ejk9q4rptzgc9djat3trtlxfvaz80g5wr6d3k",
		Amount:  math.NewInt(256700000),
		Claimed: false,
	},
	{
		Address: "elys1g5zdh7elza3f7m630fp6e8ldqxys3wjmvvu72h",
		Amount:  math.NewInt(255600000),
		Claimed: false,
	},
	{
		Address: "elys1p70p70u7h4cuq87t992sckvzff2xjxa5xmnhy3",
		Amount:  math.NewInt(255500000),
		Claimed: false,
	},
	{
		Address: "elys126lxa6kg2huh8g6w57cvzngd8k7wn40c29m38r",
		Amount:  math.NewInt(254400000),
		Claimed: false,
	},
	{
		Address: "elys1gkt3usjmfypkdnesmlja7unf3xzpxmdqjgu74e",
		Amount:  math.NewInt(253900000),
		Claimed: false,
	},
	{
		Address: "elys1frpmnepsdxmpy3amgeexguf8zqvjqkg53a5fns",
		Amount:  math.NewInt(253200000),
		Claimed: false,
	},
	{
		Address: "elys18trmku5jtmlhcx8l5hgq8wkcp2r8f9favwxkex",
		Amount:  math.NewInt(252700000),
		Claimed: false,
	},
	{
		Address: "elys1swxwy66zzel8gk6vn86rn3mqenuq89tsk8mjnd",
		Amount:  math.NewInt(252200000),
		Claimed: false,
	},
	{
		Address: "elys1ugmyqw3ac3te7969r2ppt0my545k68jz0edj84",
		Amount:  math.NewInt(251900000),
		Claimed: false,
	},
	{
		Address: "elys158em0xxgju7r090l83eley67d5z2epl7e2d6xc",
		Amount:  math.NewInt(251800000),
		Claimed: false,
	},
	{
		Address: "elys1twg04dqwd64wvafrqdal0qqc9xtcgm9kedp4u3",
		Amount:  math.NewInt(251800000),
		Claimed: false,
	},
	{
		Address: "elys1s20am2qkjcs86ueshhwc0g2zr4knqph9jhxnhk",
		Amount:  math.NewInt(251600000),
		Claimed: false,
	},
	{
		Address: "elys1w5q4q8p4u98g60y2myl54tjfpdpsf2f5w0hzz2",
		Amount:  math.NewInt(250900000),
		Claimed: false,
	},
	{
		Address: "elys1fewft0fuucrgc45ltnnl0zaygv3z5dz4qslyzc",
		Amount:  math.NewInt(250000000),
		Claimed: false,
	},
	{
		Address: "elys16mg0qrs0uj9rqefjpde24qvjelxk6sgftc3np4",
		Amount:  math.NewInt(249000000),
		Claimed: false,
	},
	{
		Address: "elys1utg24fkldlpw4nxlugqfrwzd3l32rx9glj3y2j",
		Amount:  math.NewInt(247100000),
		Claimed: false,
	},
	{
		Address: "elys194n6p9gf4zckmrzva8yrfrkgwh3egl6nystcjd",
		Amount:  math.NewInt(246600000),
		Claimed: false,
	},
	{
		Address: "elys1hyceyqw5rm0kfyca2u7zzjj4ct8gfgz5rcqmvz",
		Amount:  math.NewInt(246500000),
		Claimed: false,
	},
	{
		Address: "elys182y27eywyfejhhxnmh8gmw3h0erdd0dyja6tp4",
		Amount:  math.NewInt(245300000),
		Claimed: false,
	},
	{
		Address: "elys1y6pyrcec2pdfknz7m9l3gzsd0vzg0ahxm0l4d5",
		Amount:  math.NewInt(243900000),
		Claimed: false,
	},
	{
		Address: "elys1eh607xnhe2ewy033lqule3cpgujjpayq3khmly",
		Amount:  math.NewInt(242500000),
		Claimed: false,
	},
	{
		Address: "elys1hnrajw0md68scswszy93grycrzqrzak56urh3h",
		Amount:  math.NewInt(242400000),
		Claimed: false,
	},
	{
		Address: "elys13dzd866dtxcg3umwhmcgqpauqur5hem2rg05yv",
		Amount:  math.NewInt(242100000),
		Claimed: false,
	},
	{
		Address: "elys18km2q4dndvukf0pjluq5t2rxyrpje2ju5t2fe3",
		Amount:  math.NewInt(240500000),
		Claimed: false,
	},
	{
		Address: "elys1a06aml3dkh2z96q4vg57m43dzrv0rz89xwyln8",
		Amount:  math.NewInt(239900000),
		Claimed: false,
	},
	{
		Address: "elys1azalsnxu9z264g2cctswvwqy5tvnhf4kyq9kvu",
		Amount:  math.NewInt(239800000),
		Claimed: false,
	},
	{
		Address: "elys1x7nc6ylzn8apkrpvvs2g6mr5uvufku0w7y8eal",
		Amount:  math.NewInt(239700000),
		Claimed: false,
	},
	{
		Address: "elys1xu4qp84tyukvez9g570fajnlt43f6h8fczyga9",
		Amount:  math.NewInt(239400000),
		Claimed: false,
	},
	{
		Address: "elys1uqutvyz53epd3gqxpcr05vshce0d5p9vy6ze6d",
		Amount:  math.NewInt(238900000),
		Claimed: false,
	},
	{
		Address: "elys1egyrcy2tq4824tldn88tnhqzxl340yqlkc85ut",
		Amount:  math.NewInt(238200000),
		Claimed: false,
	},
	{
		Address: "elys1jtqa9p03ntn5tgckczv4zp0nf0vg2ej95zf25h",
		Amount:  math.NewInt(237700000),
		Claimed: false,
	},
	{
		Address: "elys1eagjx0rnqszlguweh8qkz024alxsdyldu8ejr0",
		Amount:  math.NewInt(237600000),
		Claimed: false,
	},
	{
		Address: "elys18k6qyv3kk4kx265qf6nnm4764nmmjs6xdwljnz",
		Amount:  math.NewInt(237300000),
		Claimed: false,
	},
	{
		Address: "elys19p2ya2xmx96psr876d6rp2r40e875mw3hgsha4",
		Amount:  math.NewInt(236700000),
		Claimed: false,
	},
	{
		Address: "elys1shghsu9s9u67qa4k8xl7t5e57asaj22uahnrr0",
		Amount:  math.NewInt(236300000),
		Claimed: false,
	},
	{
		Address: "elys13szrz7jkfyp4nkat36a3jes6tuqfrls98ck0pc",
		Amount:  math.NewInt(235600000),
		Claimed: false,
	},
	{
		Address: "elys1j9swqr9k5z9s8kyc4qq4vujsttgrww0x0arsxf",
		Amount:  math.NewInt(235300000),
		Claimed: false,
	},
	{
		Address: "elys1uy4q0x50tmth0wtvz0pfk62j4a2e0e7eqxuquz",
		Amount:  math.NewInt(235100000),
		Claimed: false,
	},
	{
		Address: "elys1azyaw2v5lhcl49p5czqmv9rgmagfzrcxgj5j4x",
		Amount:  math.NewInt(234900000),
		Claimed: false,
	},
	{
		Address: "elys1y4nn6p5awvz6vywzl4zl895zyftd6vqlajgqm7",
		Amount:  math.NewInt(232200000),
		Claimed: false,
	},
	{
		Address: "elys1z8l7umujphfzrnvhvnpsaa7gajffts8jvyz4mc",
		Amount:  math.NewInt(231600000),
		Claimed: false,
	},
	{
		Address: "elys1gtvaan4vngae6xafwjgkeqxwavxt3ww4yqlcme",
		Amount:  math.NewInt(230300000),
		Claimed: false,
	},
	{
		Address: "elys1cv9rf3ctkugcckt3phssd06mwc4dnuh4pw0ywa",
		Amount:  math.NewInt(230000000),
		Claimed: false,
	},
	{
		Address: "elys1vm469w820ug8cxln03wlnsnwmt2pwx828w3ajq",
		Amount:  math.NewInt(229200000),
		Claimed: false,
	},
	{
		Address: "elys12z762kemv8un7qv33mjk4pg5n67trk5lhe92xh",
		Amount:  math.NewInt(228900000),
		Claimed: false,
	},
	{
		Address: "elys1wclhrv0v8trutnjgr9u8ntjcx6fk8zarvtd9w2",
		Amount:  math.NewInt(228200000),
		Claimed: false,
	},
	{
		Address: "elys13sq0qdqt79hajetraezukj9hj2dnksctdhqgr0",
		Amount:  math.NewInt(227400000),
		Claimed: false,
	},
	{
		Address: "elys13vjz3ferqyy9up8uu3yutj0w4p0t0a6973fkuz",
		Amount:  math.NewInt(226800000),
		Claimed: false,
	},
	{
		Address: "elys1f5kn3nxyjd2vaydkczaqq40amvm7unvupr7wkt",
		Amount:  math.NewInt(226100000),
		Claimed: false,
	},
	{
		Address: "elys1y2vqwgn377urt96m5aj7hqql3trv6zjxdlzfgj",
		Amount:  math.NewInt(225500000),
		Claimed: false,
	},
	{
		Address: "elys12ystdq5d040hgyvhf3n4maj5ldzzz5lxdfam0k",
		Amount:  math.NewInt(224300000),
		Claimed: false,
	},
	{
		Address: "elys1hp7aart0up9jvsmm6w45fhxz57s4nne7w0x74r",
		Amount:  math.NewInt(224200000),
		Claimed: false,
	},
	{
		Address: "elys1dvkhzkz6vfac950lydxqfa42u8z5g8fpyjcmc7",
		Amount:  math.NewInt(222600000),
		Claimed: false,
	},
	{
		Address: "elys1ha82lhrtex4x5le4jra57du8hnw7rmw3thczvu",
		Amount:  math.NewInt(222600000),
		Claimed: false,
	},
	{
		Address: "elys12urcu3vazkqx2l7wg4p5vgcjces5xpmz09c69n",
		Amount:  math.NewInt(221800000),
		Claimed: false,
	},
	{
		Address: "elys14vvtqmsznan2y7xj5y7mm5usmc4qydpmvne48t",
		Amount:  math.NewInt(221000000),
		Claimed: false,
	},
	{
		Address: "elys1l37rhrqlqw4fdjqp322vnkt30fkcqvxmsw2yus",
		Amount:  math.NewInt(220500000),
		Claimed: false,
	},
	{
		Address: "elys1se4vef3fz2jdfqep9mzwr3nhj27qc72ggvslv8",
		Amount:  math.NewInt(220100000),
		Claimed: false,
	},
	{
		Address: "elys1st898c2x8072q9ffvyah9fxkrruqkwzmf69qa0",
		Amount:  math.NewInt(220100000),
		Claimed: false,
	},
	{
		Address: "elys1zdv659dmcewfpful2p8m995rfzfmv9pqlre68y",
		Amount:  math.NewInt(219300000),
		Claimed: false,
	},
	{
		Address: "elys1mk4ddr0v3nxgkd2e78snnd4gtuvgfeealdlkhc",
		Amount:  math.NewInt(219100000),
		Claimed: false,
	},
	{
		Address: "elys18wumeqr2hwqw9alaa0sk7s0aehalqmlgchk2qw",
		Amount:  math.NewInt(219000000),
		Claimed: false,
	},
	{
		Address: "elys1twkls4pd2ezupa02xa8ef0qrrclnkhn67w3t38",
		Amount:  math.NewInt(218300000),
		Claimed: false,
	},
	{
		Address: "elys1c278awl96ark2s9c9gz5yucrctj36pe0t597s0",
		Amount:  math.NewInt(217800000),
		Claimed: false,
	},
	{
		Address: "elys1e7wyepdpag9m8xryty464j7upw5at4a0vhfc2f",
		Amount:  math.NewInt(217600000),
		Claimed: false,
	},
	{
		Address: "elys1vwqh6aszt2vvjc2yefk8w8k3jlds4zcf38al2c",
		Amount:  math.NewInt(217600000),
		Claimed: false,
	},
	{
		Address: "elys1rwytctrqqaq08mdd4c7v57sqah44guxvm696cg",
		Amount:  math.NewInt(217500000),
		Claimed: false,
	},
	{
		Address: "elys1r9u8cj3k0gyh3aev24ynpqz2nuuepn3xtudclw",
		Amount:  math.NewInt(217300000),
		Claimed: false,
	},
	{
		Address: "elys1xtjc9gufkcdcq0m5svn4l8aguk6ktwlr8h4prg",
		Amount:  math.NewInt(216600000),
		Claimed: false,
	},
	{
		Address: "elys18taua5tzexraqx0rafaj304xvj3q8uz7r73huy",
		Amount:  math.NewInt(216500000),
		Claimed: false,
	},
	{
		Address: "elys1rww2jvedg0mhn92gpe0qn24y794hh00jxffvlm",
		Amount:  math.NewInt(216500000),
		Claimed: false,
	},
	{
		Address: "elys1mcp4s2gztp7ke0r90rfmmqknst8jmtpu8fnz4q",
		Amount:  math.NewInt(216200000),
		Claimed: false,
	},
	{
		Address: "elys1fr0vh5pja4el03spy2e0hy9mntaj0ghhtn2xam",
		Amount:  math.NewInt(215400000),
		Claimed: false,
	},
	{
		Address: "elys10nqefx0rd7he3770mjwaasd9wul6n2m9vjfcyw",
		Amount:  math.NewInt(214800000),
		Claimed: false,
	},
	{
		Address: "elys1e4x9y3atxm0g833ccfz63ua470h7fakhx0hsts",
		Amount:  math.NewInt(214600000),
		Claimed: false,
	},
	{
		Address: "elys15fkqrm3ekmtlw0aljw6jy9c62f3c2htzj4u0nr",
		Amount:  math.NewInt(214500000),
		Claimed: false,
	},
	{
		Address: "elys1dy5qcwgjtjgq6vf0ya4tl07kkr259up7k58hm7",
		Amount:  math.NewInt(214000000),
		Claimed: false,
	},
	{
		Address: "elys1ac840gsjkhupcv5qs0h3l530c84lkehkg24p2d",
		Amount:  math.NewInt(213300000),
		Claimed: false,
	},
	{
		Address: "elys12pxkw6vz523unudxk6nfmp5hlakcd28rpk67gk",
		Amount:  math.NewInt(213000000),
		Claimed: false,
	},
	{
		Address: "elys1ma62s6cwx5j72aerytu0f92ehtxff2lz3cvvr2",
		Amount:  math.NewInt(212800000),
		Claimed: false,
	},
	{
		Address: "elys17t7z6gnn9qdzh6sqr34n03hcpwkfd77urxeulx",
		Amount:  math.NewInt(211800000),
		Claimed: false,
	},
	{
		Address: "elys1k5qhpkjq8dwqfs23g7fk4x0m4yt9nlcuphxsm3",
		Amount:  math.NewInt(211400000),
		Claimed: false,
	},
	{
		Address: "elys1yqn3ywa20q3rgttszccv99smtenejge5mum7c2",
		Amount:  math.NewInt(211200000),
		Claimed: false,
	},
	{
		Address: "elys1s3n2qu9v7zquy8vus56ph06cxtgtz02rmtu7ar",
		Amount:  math.NewInt(211000000),
		Claimed: false,
	},
	{
		Address: "elys1vd6902pnx3vdxcdru89gy3rnd384vyxuflejyc",
		Amount:  math.NewInt(211000000),
		Claimed: false,
	},
	{
		Address: "elys1j0tds2clg5t0f6ams5c9c6wf3rczq8g22hpce8",
		Amount:  math.NewInt(210600000),
		Claimed: false,
	},
	{
		Address: "elys1p95nwua9ksyv6qwfen0a8jr72yyge5q0ejguqn",
		Amount:  math.NewInt(210300000),
		Claimed: false,
	},
	{
		Address: "elys14uxneahv5fuesg95s6m6ny733dvex8jl4qez6m",
		Amount:  math.NewInt(210200000),
		Claimed: false,
	},
	{
		Address: "elys1d28daguxr2yk4hztjlser99440j99f8thgvez7",
		Amount:  math.NewInt(210200000),
		Claimed: false,
	},
	{
		Address: "elys1lnje3hq9289ln32gxc9uv956l6spvr3g4qq6mj",
		Amount:  math.NewInt(209500000),
		Claimed: false,
	},
	{
		Address: "elys1sng225elftexjvwha0458663vfdpm5xuwj6xle",
		Amount:  math.NewInt(209300000),
		Claimed: false,
	},
	{
		Address: "elys164hf4ft72gemqtv5fwavzvyvn68tealuerwltc",
		Amount:  math.NewInt(208700000),
		Claimed: false,
	},
	{
		Address: "elys16l7zmlae5t5wmhxjkx3kwqy70pt05qxdsp538n",
		Amount:  math.NewInt(208600000),
		Claimed: false,
	},
	{
		Address: "elys1hdqcaue4umvkeae5ccnu6vt8ajt6l0rwpvpxmr",
		Amount:  math.NewInt(208600000),
		Claimed: false,
	},
	{
		Address: "elys1hc253sfcn66wk7884g3kh7qd9k327hxae3tq89",
		Amount:  math.NewInt(208000000),
		Claimed: false,
	},
	{
		Address: "elys170dudm95c8tavgvcl459hv3p75mzwf9k3yzg8t",
		Amount:  math.NewInt(206600000),
		Claimed: false,
	},
	{
		Address: "elys14cl93l9d6q3f7vyaddwmrt46ttlmvmvp530qqx",
		Amount:  math.NewInt(205700000),
		Claimed: false,
	},
	{
		Address: "elys1sf90h3v0lc3f6sjejscqy4nsaxj8prd43j5tkn",
		Amount:  math.NewInt(205600000),
		Claimed: false,
	},
	{
		Address: "elys1u79gq8yaj4jxewyl888q5g5xqznzp977tqdxrs",
		Amount:  math.NewInt(204100000),
		Claimed: false,
	},
	{
		Address: "elys1dlfgp4xl5s4kcqvyyzff0376c0lj8al743wxax",
		Amount:  math.NewInt(203400000),
		Claimed: false,
	},
	{
		Address: "elys1nnf85kd8d8fx4dnccrqdpr5eq6d74vhazvrpe8",
		Amount:  math.NewInt(203200000),
		Claimed: false,
	},
	{
		Address: "elys1d50f9lzduyhjkxaxnamsnsk22yxx33s9063fsk",
		Amount:  math.NewInt(202600000),
		Claimed: false,
	},
	{
		Address: "elys1kyvrk86vxnl5uh7wh6x9tw2yyf92swxg30r7fr",
		Amount:  math.NewInt(202600000),
		Claimed: false,
	},
	{
		Address: "elys1pt8atmzp05z7q32nr42zf7xu0t9gzdkydm4t9j",
		Amount:  math.NewInt(202300000),
		Claimed: false,
	},
	{
		Address: "elys1aedwh7tzxhtvh0alex86hhty5rfkwnvlc8tg24",
		Amount:  math.NewInt(202200000),
		Claimed: false,
	},
	{
		Address: "elys1hu449p59r9z2mm9uwzzejmurksavxaanhszaz2",
		Amount:  math.NewInt(201900000),
		Claimed: false,
	},
	{
		Address: "elys14mwhegj74cs90nlgdfmrgjlgd7rpq0yxcmpk5y",
		Amount:  math.NewInt(201500000),
		Claimed: false,
	},
	{
		Address: "elys1fhvnx8euurrkkskrhjkz4uhq73kex7pgskwmvj",
		Amount:  math.NewInt(201400000),
		Claimed: false,
	},
	{
		Address: "elys1v4glctg75922fvsrdzvuw0vvc72pc3g9j46swq",
		Amount:  math.NewInt(200400000),
		Claimed: false,
	},
	{
		Address: "elys1vm9eg8f0vpsfhhr73e8dzv9ulwtms4v6tp53ws",
		Amount:  math.NewInt(199800000),
		Claimed: false,
	},
	{
		Address: "elys1c4f37lt3dzrrgy6dccmcu9amnudg2aetlwmx34",
		Amount:  math.NewInt(199000000),
		Claimed: false,
	},
	{
		Address: "elys1gjr7hycp66mr550wq2tchvljnq6whjlyfku5qm",
		Amount:  math.NewInt(198300000),
		Claimed: false,
	},
	{
		Address: "elys1xwsp699f75shh9923zx9ltyu2ge26rrz94vz2m",
		Amount:  math.NewInt(198100000),
		Claimed: false,
	},
	{
		Address: "elys15nmr405vnd204tpuc8y6ghuaueqtreccvpqur6",
		Amount:  math.NewInt(197900000),
		Claimed: false,
	},
	{
		Address: "elys183nnun3xl7amzskn3apmnx7h6gwr427cunfzzj",
		Amount:  math.NewInt(197100000),
		Claimed: false,
	},
	{
		Address: "elys13c0qxfwxw2tff8sjs2hukhqpmepgj8fegrvvtf",
		Amount:  math.NewInt(196500000),
		Claimed: false,
	},
	{
		Address: "elys1anjk3uw473x0t6lrl2zf9vpmhq6xn4k9h4uz7r",
		Amount:  math.NewInt(195800000),
		Claimed: false,
	},
	{
		Address: "elys1zsyt3jx57jg06zstu8m9jj8uepecqa9x8cpkmu",
		Amount:  math.NewInt(195700000),
		Claimed: false,
	},
	{
		Address: "elys1qlljvx0ggg0qmg3ffz4f3hd0tsd259n4jd0vfx",
		Amount:  math.NewInt(195100000),
		Claimed: false,
	},
	{
		Address: "elys10maww4mrlj23atuckhljzdt3zfwu7040tz38c5",
		Amount:  math.NewInt(194000000),
		Claimed: false,
	},
	{
		Address: "elys1twlg6s7nffsrnjafyjx3mhml430z8u9erksewa",
		Amount:  math.NewInt(193900000),
		Claimed: false,
	},
	{
		Address: "elys1n0fty4he7nagltaqa0r9ad6hvc72566x82l4rj",
		Amount:  math.NewInt(193400000),
		Claimed: false,
	},
	{
		Address: "elys1hc0fqn44kngjzk8daxucmsvuq0y3rt7h6glua6",
		Amount:  math.NewInt(193000000),
		Claimed: false,
	},
	{
		Address: "elys1cn64hps3s84722vksj25atx0ngk0hy3vke3zrg",
		Amount:  math.NewInt(192800000),
		Claimed: false,
	},
	{
		Address: "elys1f5fqywq8wyssauj5d4ctu34fyjtgf5wugezulk",
		Amount:  math.NewInt(192800000),
		Claimed: false,
	},
	{
		Address: "elys1xsfh0fujn2dkwrnsak8td68yx5dqvz36802y9y",
		Amount:  math.NewInt(192800000),
		Claimed: false,
	},
	{
		Address: "elys104v6k7nlfwm5gwdepz3se5fcgflhtqg76vfa30",
		Amount:  math.NewInt(191600000),
		Claimed: false,
	},
	{
		Address: "elys1y559rcsv2vdkzmfnassmgxmfaurd3j5r4pryrx",
		Amount:  math.NewInt(191600000),
		Claimed: false,
	},
	{
		Address: "elys1ut8qdryqjlrhgfup20yuvmgtpa5892gnrn4wjk",
		Amount:  math.NewInt(190800000),
		Claimed: false,
	},
	{
		Address: "elys1vpu4qlqdgquxa454z0n3ft5ymjuftj6ktn2g9y",
		Amount:  math.NewInt(190800000),
		Claimed: false,
	},
	{
		Address: "elys1jf0rk2sh76rgzed9czk04zxwshz3j390g3g8vc",
		Amount:  math.NewInt(190100000),
		Claimed: false,
	},
	{
		Address: "elys1vyz3kc67h5hw9w2pfsfyqw3p85tm578kgzm69x",
		Amount:  math.NewInt(188900000),
		Claimed: false,
	},
	{
		Address: "elys1c0s7t35t3vc3dgm3rwus55u0hczypf9dz8n7zu",
		Amount:  math.NewInt(187900000),
		Claimed: false,
	},
	{
		Address: "elys17qg3pf85kj56z8lye25kpwvjx6cgftq0s07596",
		Amount:  math.NewInt(187700000),
		Claimed: false,
	},
	{
		Address: "elys1qzthccvl8fdm62t9843fllfntahwztm83exehe",
		Amount:  math.NewInt(187700000),
		Claimed: false,
	},
	{
		Address: "elys1007vx3gxfuh6ppx83jxwzwrnkf5u2satkhd3aq",
		Amount:  math.NewInt(187600000),
		Claimed: false,
	},
	{
		Address: "elys1sgag458e2wfsw4y2le9sa3saw0f2t3mkz5u6ff",
		Amount:  math.NewInt(187500000),
		Claimed: false,
	},
	{
		Address: "elys108t5e22lxr7z7pgdc7qaqmap5en4xfsra8pv66",
		Amount:  math.NewInt(187400000),
		Claimed: false,
	},
	{
		Address: "elys1k5hmjnga9q52ylu487pff75spc7ezw4v2j8gwl",
		Amount:  math.NewInt(187100000),
		Claimed: false,
	},
	{
		Address: "elys1h9r64kqu7pxw6r54p8yve4mlsgp9slv04urlu8",
		Amount:  math.NewInt(187000000),
		Claimed: false,
	},
	{
		Address: "elys1ja7e7wyyz6y0tyughcg7vyq9vsagwdppa04vz9",
		Amount:  math.NewInt(186900000),
		Claimed: false,
	},
	{
		Address: "elys1kalws3khhweedm49p244tts5qnqr06nuayw3yk",
		Amount:  math.NewInt(186500000),
		Claimed: false,
	},
	{
		Address: "elys1q6nmemdmd9a24rm5z6x7e0az7yp6f8lpm9lhc3",
		Amount:  math.NewInt(186200000),
		Claimed: false,
	},
	{
		Address: "elys1hhhmpjarttaqn42lqmxqxpdhnvzaaqlhzjn95g",
		Amount:  math.NewInt(186100000),
		Claimed: false,
	},
	{
		Address: "elys1f23guxp05m3d2yl53vj7thl388a4zwvvg3ahyf",
		Amount:  math.NewInt(185900000),
		Claimed: false,
	},
	{
		Address: "elys13pq2aq0g3x58wr5l96slzgjj006r70erlhntcx",
		Amount:  math.NewInt(185700000),
		Claimed: false,
	},
	{
		Address: "elys1u5nx74q6lljlndclmsh2tn865fu6rdwmewnm7g",
		Amount:  math.NewInt(185400000),
		Claimed: false,
	},
	{
		Address: "elys1zdupga9cpf35kwxnwg9ys66vyzxtrevynz3myd",
		Amount:  math.NewInt(185100000),
		Claimed: false,
	},
	{
		Address: "elys188mts2l8g67px9mq9rdge52qr55nvu5ar6eq2r",
		Amount:  math.NewInt(185000000),
		Claimed: false,
	},
	{
		Address: "elys1nqmzdegvnk4dl3kffamkc30ylnyhej9tx3v5cx",
		Amount:  math.NewInt(185000000),
		Claimed: false,
	},
	{
		Address: "elys1ld4ss9fcnze8uhzlql6pvcuuyw8k8hhf8vn8nl",
		Amount:  math.NewInt(184200000),
		Claimed: false,
	},
	{
		Address: "elys1m7k05ywuvkm6kkgmhq7ekcqw7up5fmxvtnmg3m",
		Amount:  math.NewInt(184100000),
		Claimed: false,
	},
	{
		Address: "elys1l3e33neruwmhzflq6pwqgxzyn5hwxq32ru22f3",
		Amount:  math.NewInt(183900000),
		Claimed: false,
	},
	{
		Address: "elys1lnxaqx4mdx2y88un2wlh5n38c6qva0crkq9zy0",
		Amount:  math.NewInt(183800000),
		Claimed: false,
	},
	{
		Address: "elys1lh4v8djg3wg3uylk7kg78cdyp5y2unjlquj8xk",
		Amount:  math.NewInt(183300000),
		Claimed: false,
	},
	{
		Address: "elys10zx4drkeepfe5n4hfjx7xfj7cm6c6chujfgzjz",
		Amount:  math.NewInt(183000000),
		Claimed: false,
	},
	{
		Address: "elys1k9htds6me6hq8mzdv5xsnwr4lp9nwpmt2xw906",
		Amount:  math.NewInt(182600000),
		Claimed: false,
	},
	{
		Address: "elys175vw3m9ft5y3y0daedlv4fzlag752fmrl3ux2q",
		Amount:  math.NewInt(182500000),
		Claimed: false,
	},
	{
		Address: "elys1vjndk4mkxlwjq0fkpa6gtd2q09fv76xvwhr7t7",
		Amount:  math.NewInt(182400000),
		Claimed: false,
	},
	{
		Address: "elys14arrvm2lxcukrhz3tx7k23gtvj78z2q0lnwz4s",
		Amount:  math.NewInt(182000000),
		Claimed: false,
	},
	{
		Address: "elys14smndrx7dz788dfm20kn6zp69alep8rtmx6u9x",
		Amount:  math.NewInt(182000000),
		Claimed: false,
	},
	{
		Address: "elys15l2l7fl9thtnfxmj036auz0hnfsd3txpda442c",
		Amount:  math.NewInt(181600000),
		Claimed: false,
	},
	{
		Address: "elys1e950au6qu5act07uan6q2l5m2wjpkmdma28gxh",
		Amount:  math.NewInt(181500000),
		Claimed: false,
	},
	{
		Address: "elys15kfr4d79sszqpxeak83ef0ffth758d3u24rk8g",
		Amount:  math.NewInt(181300000),
		Claimed: false,
	},
	{
		Address: "elys1dg9ryk6qt0k9uaqcu83j7rz8ynvmmnmuzn59nl",
		Amount:  math.NewInt(181000000),
		Claimed: false,
	},
	{
		Address: "elys14lde0emyuadnujwqc80hwqhjf9rvfqzzn0s5n9",
		Amount:  math.NewInt(180700000),
		Claimed: false,
	},
	{
		Address: "elys1538pvfym90svseeftvlv54zvys3c9y4hradh80",
		Amount:  math.NewInt(180500000),
		Claimed: false,
	},
	{
		Address: "elys12eyuf7nxstpfyddvq4uvh0sx3hpyzzcxuflddk",
		Amount:  math.NewInt(180400000),
		Claimed: false,
	},
	{
		Address: "elys1ttnlxu2kl3zz94canzxvcsvupvv8dyu0zgxjsh",
		Amount:  math.NewInt(180400000),
		Claimed: false,
	},
	{
		Address: "elys1e4tkauv2w2wzy449ps7lccpautt8a5fpk4gple",
		Amount:  math.NewInt(179700000),
		Claimed: false,
	},
	{
		Address: "elys1lnr5q22n6wc4a3t6c5lldvc5737tjlpj4w8s5f",
		Amount:  math.NewInt(179700000),
		Claimed: false,
	},
	{
		Address: "elys1q95x2lwst53p292axadj5vdfzrlgvcj4dneya9",
		Amount:  math.NewInt(179400000),
		Claimed: false,
	},
	{
		Address: "elys10us6k0d9c8dc3lnssa0gwlsk6sumya6appkpea",
		Amount:  math.NewInt(179300000),
		Claimed: false,
	},
	{
		Address: "elys1lm77mf5fhz5lyl6q466fk9847efuy456qxla0m",
		Amount:  math.NewInt(179300000),
		Claimed: false,
	},
	{
		Address: "elys1vu05t9ksa46qhv2l69t8gcejnfrq70vx4vk40h",
		Amount:  math.NewInt(179300000),
		Claimed: false,
	},
	{
		Address: "elys1nw97tyt5eslcsexev2wk5x9fqnth92tv26a57z",
		Amount:  math.NewInt(178800000),
		Claimed: false,
	},
	{
		Address: "elys1pg6f5w8hr2mqhlda08tltr5zret0qdrjp3zng6",
		Amount:  math.NewInt(178400000),
		Claimed: false,
	},
	{
		Address: "elys18pj4k54r2d55cgldp5nsrh85dma5pf97myc5hp",
		Amount:  math.NewInt(178200000),
		Claimed: false,
	},
	{
		Address: "elys127665wuj3lyqxfjwnw7hfyzfmjzd7fvqs4l4dv",
		Amount:  math.NewInt(177700000),
		Claimed: false,
	},
	{
		Address: "elys1dy9eq3sfp98q34zuvrculsa3curjdyw08kv73f",
		Amount:  math.NewInt(177400000),
		Claimed: false,
	},
	{
		Address: "elys1wwcwnnea5s58m8aennlc7d6m0d9hr0g7tcgnc5",
		Amount:  math.NewInt(177300000),
		Claimed: false,
	},
	{
		Address: "elys1kxexcfqxn9ff8kqhzh7857tdpsmw0ajh84qfef",
		Amount:  math.NewInt(176500000),
		Claimed: false,
	},
	{
		Address: "elys1j2zddtrk3jypqpt8h0qh9prf0wg3uyy0yplhp2",
		Amount:  math.NewInt(176400000),
		Claimed: false,
	},
	{
		Address: "elys1pyuwlv0r30ppys865dpq7u0v46gm90u56t5jhj",
		Amount:  math.NewInt(176400000),
		Claimed: false,
	},
	{
		Address: "elys1gg7423yff8mqae6fj73vpxv74q5dgzzeqq407d",
		Amount:  math.NewInt(175900000),
		Claimed: false,
	},
	{
		Address: "elys1py5dax4v54ngscnz4fxz8jwu4at06dnnmmadpm",
		Amount:  math.NewInt(175300000),
		Claimed: false,
	},
	{
		Address: "elys1qnm7vjlp06set0r9dxmlxlssz4l7l9qmr6t2v6",
		Amount:  math.NewInt(174900000),
		Claimed: false,
	},
	{
		Address: "elys1fcyec6nkyzjxle86mcwp7ueanuy38v3dl0x8mt",
		Amount:  math.NewInt(174700000),
		Claimed: false,
	},
	{
		Address: "elys1sqm0rm2w52lg07g0v9kynt6zqqghlfvu9knyz4",
		Amount:  math.NewInt(174700000),
		Claimed: false,
	},
	{
		Address: "elys1xfj2747qa5cwy2zqa4dxm7kwpavv20mt230tyg",
		Amount:  math.NewInt(173900000),
		Claimed: false,
	},
	{
		Address: "elys1remj0f34kdrrsd5932aw87da6x8fxua677k8tv",
		Amount:  math.NewInt(173800000),
		Claimed: false,
	},
	{
		Address: "elys1pk5r2ygdegaj4m9yrhxllsketw3s0hszwk0umf",
		Amount:  math.NewInt(173700000),
		Claimed: false,
	},
	{
		Address: "elys10dp29jukrsvgsnac9my784f25hqceh59alnvef",
		Amount:  math.NewInt(173100000),
		Claimed: false,
	},
	{
		Address: "elys1c2carfuccpm62jnsrx9y7d8maj8dmy5rxpcs7y",
		Amount:  math.NewInt(173100000),
		Claimed: false,
	},
	{
		Address: "elys14ke7nqp32z3je2wkrhxp8nvd6hznq6f4lg9lzg",
		Amount:  math.NewInt(172400000),
		Claimed: false,
	},
	{
		Address: "elys1k5dwx37dj2fyk7lg72aspdy9x0tyxzyn3qj0ng",
		Amount:  math.NewInt(171700000),
		Claimed: false,
	},
	{
		Address: "elys1r9ua3ulh34d8xpzql535zyse84vtn430wpcu85",
		Amount:  math.NewInt(171700000),
		Claimed: false,
	},
	{
		Address: "elys1x0nwdy980ew5jkz2dwglud3jldn0346sytc2j3",
		Amount:  math.NewInt(171700000),
		Claimed: false,
	},
	{
		Address: "elys1eshcujwa86dqjsta2yq2gq7wemjs28jdnv5pm6",
		Amount:  math.NewInt(171500000),
		Claimed: false,
	},
	{
		Address: "elys1f7mvtw5383yyuwddh08u08cy7j49qggw8j2t2y",
		Amount:  math.NewInt(171300000),
		Claimed: false,
	},
	{
		Address: "elys13wsuwkj4xhz5skfzpczq8jmzfgqgl9vkxn8r8r",
		Amount:  math.NewInt(171000000),
		Claimed: false,
	},
	{
		Address: "elys10gwdzg8vxl6cr5dlfswe8qatkxdfxaue43eq63",
		Amount:  math.NewInt(170000000),
		Claimed: false,
	},
	{
		Address: "elys1h86zg7x940yu0ddsmr8fr4fz3c3paf3fjp96zr",
		Amount:  math.NewInt(169700000),
		Claimed: false,
	},
	{
		Address: "elys1kegaa93gqugelg4zmrrcfnjqnt4v6entv67swq",
		Amount:  math.NewInt(169600000),
		Claimed: false,
	},
	{
		Address: "elys1jhfww5dcsnjvtujs4fsfn0c9hj8lxhdj5tckap",
		Amount:  math.NewInt(169000000),
		Claimed: false,
	},
	{
		Address: "elys1fkl4tzf6vdv6rqjusgznl5upads8hwmmtrdfuf",
		Amount:  math.NewInt(168800000),
		Claimed: false,
	},
	{
		Address: "elys14wg206u3f49lx0pc5y7ns5am2hgn67cegcvghf",
		Amount:  math.NewInt(168600000),
		Claimed: false,
	},
	{
		Address: "elys1cudnvrfr8frayurauv2vh0y8ewnqvhxdwp9pss",
		Amount:  math.NewInt(168400000),
		Claimed: false,
	},
	{
		Address: "elys18w6pkfsvktjnccysylcx48lyxnfzp33qkkfv65",
		Amount:  math.NewInt(167600000),
		Claimed: false,
	},
	{
		Address: "elys1l8eaaythmuesmlhqhpzuee0ta0umkyxfam0uzu",
		Amount:  math.NewInt(167500000),
		Claimed: false,
	},
	{
		Address: "elys1hn4elntxqcpsh23ykerp5mrutn9zfq9vg9t22x",
		Amount:  math.NewInt(167300000),
		Claimed: false,
	},
	{
		Address: "elys1m88xftx7292nsg7dskxv7feckysn6cjuw43v38",
		Amount:  math.NewInt(167300000),
		Claimed: false,
	},
	{
		Address: "elys10zqsa2r4zlwdqazxe97dmwaz23pcwv4cswrqey",
		Amount:  math.NewInt(167200000),
		Claimed: false,
	},
	{
		Address: "elys1kxkp0mhk2m6mqu2jd8gqegzg9us6nlnssy3drn",
		Amount:  math.NewInt(167200000),
		Claimed: false,
	},
	{
		Address: "elys17gzhl0puuffu2s3h055pugrrn2cgqezw4kzg2t",
		Amount:  math.NewInt(166700000),
		Claimed: false,
	},
	{
		Address: "elys1lskktq94m7q79xxfpkps8fclxhsx0thf2ewckp",
		Amount:  math.NewInt(166600000),
		Claimed: false,
	},
	{
		Address: "elys18gssqdzkm3wcacls5u8ajlzht7ukxgctglm3ha",
		Amount:  math.NewInt(166300000),
		Claimed: false,
	},
	{
		Address: "elys1xjvp6skc33tm3snk7e640fzezhnajvnvddfwa3",
		Amount:  math.NewInt(166000000),
		Claimed: false,
	},
	{
		Address: "elys1u76ktndzaa2ke43rgc44yhlxntkrzdn95zcmtl",
		Amount:  math.NewInt(165900000),
		Claimed: false,
	},
	{
		Address: "elys1ql7apde38p938g93cpaqt6043ners28z4ds300",
		Amount:  math.NewInt(165400000),
		Claimed: false,
	},
	{
		Address: "elys1k565wns5qmv2ekt6tv3flh3hwquwdedawcda3a",
		Amount:  math.NewInt(165200000),
		Claimed: false,
	},
	{
		Address: "elys1ja86k72hx88v93pryjesh7g4572z5pwgkgjzfr",
		Amount:  math.NewInt(165100000),
		Claimed: false,
	},
	{
		Address: "elys1ea70t23rmxtq3uvejztgzgdu4l6atvnv36fhs3",
		Amount:  math.NewInt(164900000),
		Claimed: false,
	},
	{
		Address: "elys139yufsrqszv0ne6et009v4lmrz8p6ypsnltsaq",
		Amount:  math.NewInt(164500000),
		Claimed: false,
	},
	{
		Address: "elys1jj4649u35pqj7dt50jct28wun46jkp5crla3fe",
		Amount:  math.NewInt(164100000),
		Claimed: false,
	},
	{
		Address: "elys183kkpe4gjnpmv029empnsf4qlr3mfguljaa02v",
		Amount:  math.NewInt(163800000),
		Claimed: false,
	},
	{
		Address: "elys1c0k4cd2yvqh0c3c3p5j6dahq4gxz35tc5utwmg",
		Amount:  math.NewInt(163700000),
		Claimed: false,
	},
	{
		Address: "elys18apwxna5v04f6wmhkdwcy2mz2l3zwhgl5xqr2j",
		Amount:  math.NewInt(163600000),
		Claimed: false,
	},
	{
		Address: "elys132cd92yp6vmmqrt6vsp02yglecfjuk0hhz8hjw",
		Amount:  math.NewInt(163500000),
		Claimed: false,
	},
	{
		Address: "elys1e0zvjga2aznt6qx9aukjdmvua49dadeap2ldtf",
		Amount:  math.NewInt(163500000),
		Claimed: false,
	},
	{
		Address: "elys1strsrwawy66vpflvmgcnttep6tzz20qz238xj9",
		Amount:  math.NewInt(163300000),
		Claimed: false,
	},
	{
		Address: "elys19vecg2np7009syhqxpxqyuhk8z0lz7z5hprr4d",
		Amount:  math.NewInt(162800000),
		Claimed: false,
	},
	{
		Address: "elys1hv9rk9ltj5hq44yl4v356mlhcawpyt7s6jcpdw",
		Amount:  math.NewInt(162600000),
		Claimed: false,
	},
	{
		Address: "elys1qkfdrmekn5thhv5w7s58fsfwd306k89nhjxmcv",
		Amount:  math.NewInt(162500000),
		Claimed: false,
	},
	{
		Address: "elys1rj2c99vmxqq8drr2ns5lsmthedsgwfc43mrapd",
		Amount:  math.NewInt(162400000),
		Claimed: false,
	},
	{
		Address: "elys1g9nnxvzfuy6ak9j9kkjlxhvn39xxngqdq7y64m",
		Amount:  math.NewInt(162200000),
		Claimed: false,
	},
	{
		Address: "elys14qzzcnnyuxxl056ms757y08jzekzpjlec7ln9l",
		Amount:  math.NewInt(162100000),
		Claimed: false,
	},
	{
		Address: "elys1tw6rspjakn9wnjvlkdasw6w57gq396csde6z3f",
		Amount:  math.NewInt(161900000),
		Claimed: false,
	},
	{
		Address: "elys1synn5wch42vuxmfsfmtjz65yldn6ukzfkjg2nh",
		Amount:  math.NewInt(161800000),
		Claimed: false,
	},
	{
		Address: "elys1yyk9d4hyfs62pjkps62am54tjwz45zmku5h3xg",
		Amount:  math.NewInt(161500000),
		Claimed: false,
	},
	{
		Address: "elys1tw9cha2f434c4nfzx2euxxlcdx0832zav4f5ht",
		Amount:  math.NewInt(161400000),
		Claimed: false,
	},
	{
		Address: "elys1ap9pqr0yl762uw7tcn26pfs48wtv57d3txedcv",
		Amount:  math.NewInt(161300000),
		Claimed: false,
	},
	{
		Address: "elys1gq4zj5wsfyucp29kdjdlar9zgrga9pqykvf3us",
		Amount:  math.NewInt(161300000),
		Claimed: false,
	},
	{
		Address: "elys1j04u3kzzz77fl3yvpvcuzftqz5g0d0smjmltmx",
		Amount:  math.NewInt(161300000),
		Claimed: false,
	},
	{
		Address: "elys1j9gjtjv9x5lszh7ptl0feqff4l5xk09xzsuszg",
		Amount:  math.NewInt(161300000),
		Claimed: false,
	},
	{
		Address: "elys1euhjqf276n4cskwwejadztsj6aezn2stsek0kd",
		Amount:  math.NewInt(161100000),
		Claimed: false,
	},
	{
		Address: "elys13vxywzzlknflkj0udg4xr6wd30fjhvnsnel9zv",
		Amount:  math.NewInt(160400000),
		Claimed: false,
	},
	{
		Address: "elys1v0xjkkchxv2su584tjm89mxvmsfh2yg3qkrkwf",
		Amount:  math.NewInt(160400000),
		Claimed: false,
	},
	{
		Address: "elys1ry4zth04s5kgzg9tf0ru8gs6c9rpryj0errfja",
		Amount:  math.NewInt(160000000),
		Claimed: false,
	},
	{
		Address: "elys1gcnsk9p5wgstdafj88c3k2ty6yfvak73nnap62",
		Amount:  math.NewInt(159600000),
		Claimed: false,
	},
	{
		Address: "elys1fudxzahc0g563sentcptgdt5pdrwehcvqwr7cj",
		Amount:  math.NewInt(158800000),
		Claimed: false,
	},
	{
		Address: "elys1cwu8narxpfaldvgr4dnpakgdey7s5u87reth9j",
		Amount:  math.NewInt(157700000),
		Claimed: false,
	},
	{
		Address: "elys1x4lqa4c08ldeaa3v487pz2gwsah3c57wqyjd4s",
		Amount:  math.NewInt(157700000),
		Claimed: false,
	},
	{
		Address: "elys16ymdcn3k4frzx93d3s299v0gez5aeq6j6j3qul",
		Amount:  math.NewInt(157500000),
		Claimed: false,
	},
	{
		Address: "elys1hjpnramx58gctacs0064w6xda37nkyaugdk6pc",
		Amount:  math.NewInt(157300000),
		Claimed: false,
	},
	{
		Address: "elys1fc3m03ssm5t5ejxlxfqvvl9ttzghu2qcrmckul",
		Amount:  math.NewInt(157200000),
		Claimed: false,
	},
	{
		Address: "elys100jr7pku3at90pw9hcfep9cup3rpddkmxhav23",
		Amount:  math.NewInt(157000000),
		Claimed: false,
	},
	{
		Address: "elys15x9jgefpchudz7c3dha8mculp823hfqys2zf3v",
		Amount:  math.NewInt(157000000),
		Claimed: false,
	},
	{
		Address: "elys103q8qurph4mylv3znzcmpl396a06cplfrrufgy",
		Amount:  math.NewInt(156700000),
		Claimed: false,
	},
	{
		Address: "elys1q9kkh0wj56mckpu0m2skfkws8leeg2du2798hf",
		Amount:  math.NewInt(156600000),
		Claimed: false,
	},
	{
		Address: "elys1azp99943wtc7yuml7vcu7p9u5fanq4cvsk3c2c",
		Amount:  math.NewInt(155900000),
		Claimed: false,
	},
	{
		Address: "elys1w5kxwl9nvxjy5gf5jcm7xpff9dq0mf9530kjr8",
		Amount:  math.NewInt(155600000),
		Claimed: false,
	},
	{
		Address: "elys1057xrjpae5pk8xqjywjjr20z8r6f8a8s3c5lsq",
		Amount:  math.NewInt(155500000),
		Claimed: false,
	},
	{
		Address: "elys1al97hw8jtddyvln99j6qxzyfx8aur83zadzk5q",
		Amount:  math.NewInt(155500000),
		Claimed: false,
	},
	{
		Address: "elys1pd5cgu3rhu8xu3ulkaalzknj4c03jszvp0w29g",
		Amount:  math.NewInt(155100000),
		Claimed: false,
	},
	{
		Address: "elys1edxvcc5mg4y66ljps3lanj6lu7h7pfj5fnds6r",
		Amount:  math.NewInt(155000000),
		Claimed: false,
	},
	{
		Address: "elys16nhl3gxl7nu8p5ktlp8epkptth9t6sayjgvgya",
		Amount:  math.NewInt(154800000),
		Claimed: false,
	},
	{
		Address: "elys1rwcqhextdj3ngsyj0azgth3nl7xv32lt2vkula",
		Amount:  math.NewInt(154800000),
		Claimed: false,
	},
	{
		Address: "elys12dctfqxf2mqv0prq6n5wnkjm8y4rnlfmnu4p98",
		Amount:  math.NewInt(154600000),
		Claimed: false,
	},
	{
		Address: "elys1cr2g5s9j72ranr79z3mqxgwk8mwwjegu74sdrm",
		Amount:  math.NewInt(154500000),
		Claimed: false,
	},
	{
		Address: "elys1gk4ss97lzs672wwe72eh6y6jrh56r8lkk3e4jf",
		Amount:  math.NewInt(154100000),
		Claimed: false,
	},
	{
		Address: "elys1j04w5pmfyxfdc8wtma4jsmdzxgpd3kkpra5hty",
		Amount:  math.NewInt(154000000),
		Claimed: false,
	},
	{
		Address: "elys1n9cdavlkp6jnag0328m6c2zfsedwtsc58y932j",
		Amount:  math.NewInt(153800000),
		Claimed: false,
	},
	{
		Address: "elys1uu9q0lcjrse7en59x5m88c94ef070sh8z2fu5t",
		Amount:  math.NewInt(153800000),
		Claimed: false,
	},
	{
		Address: "elys1047xjx7x35g7myuf5cejq006nz22c5lll6zkj8",
		Amount:  math.NewInt(153600000),
		Claimed: false,
	},
	{
		Address: "elys13e9809rfjf44ah05h09lncrr9yhw2xc79vg4xx",
		Amount:  math.NewInt(153400000),
		Claimed: false,
	},
	{
		Address: "elys17qzz6ph06tgnf602378ejuy5pxzmmh7n4l6zdd",
		Amount:  math.NewInt(153100000),
		Claimed: false,
	},
	{
		Address: "elys1f3xnt2t6qtqlwcqu42elxncyenqglu9f5tkrhw",
		Amount:  math.NewInt(152700000),
		Claimed: false,
	},
	{
		Address: "elys1wv2yr3lv9nffqpy8a30lxuc2qw8lhq4a3w90zm",
		Amount:  math.NewInt(152600000),
		Claimed: false,
	},
	{
		Address: "elys16v3f9t68e78jcqn7lyuaqm5ssep56eut0gm2up",
		Amount:  math.NewInt(152400000),
		Claimed: false,
	},
	{
		Address: "elys155u3gavfezgsnucetmyam3lsnwx243d0ew5t0y",
		Amount:  math.NewInt(152100000),
		Claimed: false,
	},
	{
		Address: "elys10t3n8xz3rr88u9y4fxpqyswutt5w6kdx5fmslq",
		Amount:  math.NewInt(151700000),
		Claimed: false,
	},
	{
		Address: "elys1zerql52xeuqhslcyet0ek2g7w3nav9m7nnctdl",
		Amount:  math.NewInt(151700000),
		Claimed: false,
	},
	{
		Address: "elys1vmc4zavdt07ggq4v6mwgycpr22xl0fd64za6vp",
		Amount:  math.NewInt(151300000),
		Claimed: false,
	},
	{
		Address: "elys12mjrcz9ftclzut9g4ruwu9hm24ymq5d44q0gyy",
		Amount:  math.NewInt(150400000),
		Claimed: false,
	},
	{
		Address: "elys1cj8wah7vngfh2wpee8558u94uy9n7zmenhwp34",
		Amount:  math.NewInt(150400000),
		Claimed: false,
	},
	{
		Address: "elys15yed837vffpgkmxqmmpvhm8hezq9xxzpe687nd",
		Amount:  math.NewInt(149700000),
		Claimed: false,
	},
	{
		Address: "elys1j40vwe46c2x7wlnuj0hzkgwyh0xjygsnhwqvh0",
		Amount:  math.NewInt(149700000),
		Claimed: false,
	},
	{
		Address: "elys194mn4l28dv2xu28h4es89wmx04h4zy40hy9q4s",
		Amount:  math.NewInt(149400000),
		Claimed: false,
	},
	{
		Address: "elys1ksa8h90sa0wwmea8ms4vvz4xdzdttave6hz7ks",
		Amount:  math.NewInt(149400000),
		Claimed: false,
	},
	{
		Address: "elys14ayev8lq0r9ngmmcrky4sh9f9xl4qeg6ax0swy",
		Amount:  math.NewInt(149200000),
		Claimed: false,
	},
	{
		Address: "elys17qacjuwz58vupckhyz0g20f742khmfxs57yvjp",
		Amount:  math.NewInt(149200000),
		Claimed: false,
	},
	{
		Address: "elys1s2cg9a9yxwrljjmyulvsmshv70k0dlne0etjsm",
		Amount:  math.NewInt(149200000),
		Claimed: false,
	},
	{
		Address: "elys1gh7m789a03rspsvr94nj3747ysju85u5g0ndps",
		Amount:  math.NewInt(149100000),
		Claimed: false,
	},
	{
		Address: "elys1ty9p2d8z36ancxl9hflyqzkxnc06x2pfxhuawu",
		Amount:  math.NewInt(149000000),
		Claimed: false,
	},
	{
		Address: "elys10vct8xc6rvgddzyk4pwh590gm25ttkdvxr3fch",
		Amount:  math.NewInt(148900000),
		Claimed: false,
	},
	{
		Address: "elys12t9w6gjyq54zms3hm7wg0034nrvk7yzcdcqqsr",
		Amount:  math.NewInt(148000000),
		Claimed: false,
	},
	{
		Address: "elys1v8hu36fsu7asarfsy5sa5k64c4qcldvm5sdj5d",
		Amount:  math.NewInt(147700000),
		Claimed: false,
	},
	{
		Address: "elys1w69h3k5vfg6wngjh882tc9q6ymq46695qp6mq7",
		Amount:  math.NewInt(147700000),
		Claimed: false,
	},
	{
		Address: "elys1w8th0yh7asefv0wv44s4xd5tn7wwck0ptuaa7e",
		Amount:  math.NewInt(147500000),
		Claimed: false,
	},
	{
		Address: "elys1tkaw3mwkp5f2jx07jmxsu6srr5shsh7sqats3y",
		Amount:  math.NewInt(147300000),
		Claimed: false,
	},
	{
		Address: "elys1u4338fg3tqw32fqntnkrfpwqg25s2s57y3ltvd",
		Amount:  math.NewInt(147200000),
		Claimed: false,
	},
	{
		Address: "elys1q0lv6r62ralqu5fsf6mm3lax3v25jm2el6ra0y",
		Amount:  math.NewInt(147100000),
		Claimed: false,
	},
	{
		Address: "elys1xf6a86dqavtcg3dcuulu2u6khsezjuaxnmuh60",
		Amount:  math.NewInt(146900000),
		Claimed: false,
	},
	{
		Address: "elys1yqwn0u83mv3fttknk6krkrd9d0pz4dx6pptrpm",
		Amount:  math.NewInt(146700000),
		Claimed: false,
	},
	{
		Address: "elys1k5lhcpr0fmu8wg6skarhsdvdj50fjjwqg60xyy",
		Amount:  math.NewInt(146500000),
		Claimed: false,
	},
	{
		Address: "elys1rvavftd6apy234nde6t0rnrf49savphq8vy2ek",
		Amount:  math.NewInt(146400000),
		Claimed: false,
	},
	{
		Address: "elys165ahz8ys420mqn4335eca2puf5pj4x5g57lm6l",
		Amount:  math.NewInt(146300000),
		Claimed: false,
	},
	{
		Address: "elys1rxzj2xsvj9zspyuwkhvuwwajv338gzlcs6023q",
		Amount:  math.NewInt(146300000),
		Claimed: false,
	},
	{
		Address: "elys14jk06dv6hhrmvjcrxt504m90kkrzmj6mnw3j75",
		Amount:  math.NewInt(146000000),
		Claimed: false,
	},
	{
		Address: "elys14vpgkpf7r9c4ymgmq9w9rzf5setvk23wzxs5sk",
		Amount:  math.NewInt(146000000),
		Claimed: false,
	},
	{
		Address: "elys1x4l0gv8gvgqhh83qguahrfstazsavc4p29xtgs",
		Amount:  math.NewInt(145900000),
		Claimed: false,
	},
	{
		Address: "elys1sxl39y0t409slktjsfzt2z4qz7yu0pl205pn3c",
		Amount:  math.NewInt(145500000),
		Claimed: false,
	},
	{
		Address: "elys12ypr6398q7urddnzt25gzc0xf7zhs4yuwdnxtm",
		Amount:  math.NewInt(145400000),
		Claimed: false,
	},
	{
		Address: "elys1p4eh7hpd3ejt7a3pjz7mq4j0fjvf3jpv8lkjtj",
		Amount:  math.NewInt(144900000),
		Claimed: false,
	},
	{
		Address: "elys1yrcc9tyqu6sl9q4knvtrt2u9hhp32a67ua628m",
		Amount:  math.NewInt(144700000),
		Claimed: false,
	},
	{
		Address: "elys1srpe4ka9w4vtlszt3lj5rn0x5gjzfg9xf5r3ck",
		Amount:  math.NewInt(144500000),
		Claimed: false,
	},
	{
		Address: "elys1mfpccfvyvf2qe3x692dpywrpnaajzmryuv7m5e",
		Amount:  math.NewInt(144300000),
		Claimed: false,
	},
	{
		Address: "elys18mgh8uregz9785uqapwrp8ke6663fyfqds2ewn",
		Amount:  math.NewInt(144200000),
		Claimed: false,
	},
	{
		Address: "elys1nq32l4drp2t28ntlq95n9y65kzry6sxzjnrcl9",
		Amount:  math.NewInt(144200000),
		Claimed: false,
	},
	{
		Address: "elys1rlclqgc6r5cj7gnlr62nke59xre25zn2etut4r",
		Amount:  math.NewInt(143900000),
		Claimed: false,
	},
	{
		Address: "elys15e48k6j8hsvlvpy58qjtdprwpypankx4kf4yex",
		Amount:  math.NewInt(143700000),
		Claimed: false,
	},
	{
		Address: "elys1v9d8rvgppejj9vufuflf59sx542nuanry5ma9g",
		Amount:  math.NewInt(143700000),
		Claimed: false,
	},
	{
		Address: "elys1vujuzzps6m5tcjc5r84yum5gh9g5vys38smnsn",
		Amount:  math.NewInt(143700000),
		Claimed: false,
	},
	{
		Address: "elys12djnh2za7zyxmeh87m0qsvkhkqyzwml70ywx9u",
		Amount:  math.NewInt(143600000),
		Claimed: false,
	},
	{
		Address: "elys1434zluxmfj85vwpc3stk03wrcuzrqllc88ztz0",
		Amount:  math.NewInt(143600000),
		Claimed: false,
	},
	{
		Address: "elys19pj4suq9xjsnrus9pv285kw674uckvfnsvqlgk",
		Amount:  math.NewInt(143600000),
		Claimed: false,
	},
	{
		Address: "elys1d98g9wm2cgduxfct40aavpgczvzpw85zlwmddd",
		Amount:  math.NewInt(143600000),
		Claimed: false,
	},
	{
		Address: "elys1ez4ev32pznjc6lp7yakdc7yj2quqhjrf4m4gqr",
		Amount:  math.NewInt(143600000),
		Claimed: false,
	},
	{
		Address: "elys1fmmkme8sa5a9lvf5z4dqv8slk7je26cpdzx0q6",
		Amount:  math.NewInt(143600000),
		Claimed: false,
	},
	{
		Address: "elys1k67qnlfehmxcc40yqntc864q7c94jd8gt6laq9",
		Amount:  math.NewInt(143600000),
		Claimed: false,
	},
	{
		Address: "elys1lh6ruf253fhjykdesvfccy0y0vhjque4gjkrzz",
		Amount:  math.NewInt(143600000),
		Claimed: false,
	},
	{
		Address: "elys1r0mlu77nlkqlll55uvqfpsnyyt9ndcrltwgpjq",
		Amount:  math.NewInt(143600000),
		Claimed: false,
	},
	{
		Address: "elys1ujw7fswvr2mcaz72zchsdljwk3rxk6xgjqs9hf",
		Amount:  math.NewInt(143600000),
		Claimed: false,
	},
	{
		Address: "elys1gytmmztflnm3naxwwst0jh4lnpu5lw9z0y309u",
		Amount:  math.NewInt(143500000),
		Claimed: false,
	},
	{
		Address: "elys15kclc9asxfrqly45d4xcalv2xlsxehfzla475j",
		Amount:  math.NewInt(143200000),
		Claimed: false,
	},
	{
		Address: "elys1al527kzge9yf2ld5em8ln5e3xzd9jgp94zpnwr",
		Amount:  math.NewInt(143200000),
		Claimed: false,
	},
	{
		Address: "elys1awfqa6xjgzcfe9w5aaaz7trjcpvd92fjtcv824",
		Amount:  math.NewInt(143200000),
		Claimed: false,
	},
	{
		Address: "elys14yeut6j3weapa0nwjg377xs2d6ac7gk6l5y0x2",
		Amount:  math.NewInt(142900000),
		Claimed: false,
	},
	{
		Address: "elys1nw33tg8j5hhxdszhucw9s0fpm4wyrmddz0wzx9",
		Amount:  math.NewInt(142900000),
		Claimed: false,
	},
	{
		Address: "elys10nr605yasthfc505q83z5tv5lyshwl9gv36gsd",
		Amount:  math.NewInt(142800000),
		Claimed: false,
	},
	{
		Address: "elys17v6apkcxlv4qsrh8vgnvlccwy7mjyzkyg3k5a5",
		Amount:  math.NewInt(142500000),
		Claimed: false,
	},
	{
		Address: "elys1cmjstdv8ye22fwpq3780hkwxjl4zuqual9qjuu",
		Amount:  math.NewInt(142500000),
		Claimed: false,
	},
	{
		Address: "elys1f2lqzwx5h08xx75jyqf608cquyj6ajud80kpc6",
		Amount:  math.NewInt(142500000),
		Claimed: false,
	},
	{
		Address: "elys1fcsuxtwyem4fzc4ta7scgz5u5x87h2xdd3ygkh",
		Amount:  math.NewInt(142500000),
		Claimed: false,
	},
	{
		Address: "elys1n2nx2kagpxt6ge92vhdufnh9jyvqkfdj444gx9",
		Amount:  math.NewInt(142500000),
		Claimed: false,
	},
	{
		Address: "elys1n6a2xrcrnsvru8pwgtv4537kamsdf6yed3d54v",
		Amount:  math.NewInt(142500000),
		Claimed: false,
	},
	{
		Address: "elys1ngy08w0k6u02gse54fvqt3mzs8y63hxvsxp574",
		Amount:  math.NewInt(142500000),
		Claimed: false,
	},
	{
		Address: "elys1unqfc8zm3cl7cf0zv9zn3rsrrsrju6f8mkyz0t",
		Amount:  math.NewInt(142500000),
		Claimed: false,
	},
	{
		Address: "elys1x4f0s9p7gha9rjr4q2n2f592u5vp0cuyv2vrgy",
		Amount:  math.NewInt(142500000),
		Claimed: false,
	},
	{
		Address: "elys1zcy7h9vpj2g4w0xw3ta3ky977z6e6ht093cgwz",
		Amount:  math.NewInt(142500000),
		Claimed: false,
	},
	{
		Address: "elys1688nuyz9q9cjmj83fh8dy8wesk7lp8czvph3lh",
		Amount:  math.NewInt(142400000),
		Claimed: false,
	},
	{
		Address: "elys184kccexnk3gc0d6llsr7pxvlkcmy34yfks6wa8",
		Amount:  math.NewInt(142400000),
		Claimed: false,
	},
	{
		Address: "elys19edwjr6ypnmqcl4kdmx36qvqa0d5xdx70xh0ec",
		Amount:  math.NewInt(142300000),
		Claimed: false,
	},
	{
		Address: "elys1mc4frfqrvna3x6u279xerz37pc4pnyxm7ukggn",
		Amount:  math.NewInt(142300000),
		Claimed: false,
	},
	{
		Address: "elys12kje68dhtumzhv7mmd8s0pc50wsn00627sx9mt",
		Amount:  math.NewInt(142200000),
		Claimed: false,
	},
	{
		Address: "elys1djjzq8expg9t8v9sk5wz59guhw2ewwv9gejzw0",
		Amount:  math.NewInt(142200000),
		Claimed: false,
	},
	{
		Address: "elys19na22kzzlhrmny70f5mcg33dhzqz79w3650muv",
		Amount:  math.NewInt(141700000),
		Claimed: false,
	},
	{
		Address: "elys1etvj49upnsg45hyds0qqmtv9x49s86vmajl89g",
		Amount:  math.NewInt(141600000),
		Claimed: false,
	},
	{
		Address: "elys1uc0jl4gq6uuf344s0ppsynqh0hnch4cmadzw35",
		Amount:  math.NewInt(141600000),
		Claimed: false,
	},
	{
		Address: "elys16c2930el0gd83rw89wmq9nadmr27tl5rzqvw6s",
		Amount:  math.NewInt(141300000),
		Claimed: false,
	},
	{
		Address: "elys1v7a0z06sxwlmyqekn4uq79d0l3v08zdlrpw5ju",
		Amount:  math.NewInt(141300000),
		Claimed: false,
	},
	{
		Address: "elys1p69wkfuxavnxwkdh27c5uxvx70m3j3j2mnepgv",
		Amount:  math.NewInt(141200000),
		Claimed: false,
	},
	{
		Address: "elys1488kv08va80uyw3rqm03dxapmgak6mx47kl5hz",
		Amount:  math.NewInt(141100000),
		Claimed: false,
	},
	{
		Address: "elys167ggfanrwmzqj3j3hv5pmnvc478a0rvud8mnds",
		Amount:  math.NewInt(141000000),
		Claimed: false,
	},
	{
		Address: "elys1a3gxplhje6f9k0605x5qwhd9cwwpjnvsvd0krg",
		Amount:  math.NewInt(140800000),
		Claimed: false,
	},
	{
		Address: "elys197x574luwjh5ngp43www2t7xjhvu4xr4q8zdlf",
		Amount:  math.NewInt(140600000),
		Claimed: false,
	},
	{
		Address: "elys1yeh00zd5zjts4rx2r0w0e4vxpz9c59dfxedtdr",
		Amount:  math.NewInt(140600000),
		Claimed: false,
	},
	{
		Address: "elys1zg6dr6rs90ad3t0k9y3sr6kq6j0wdf358jluyj",
		Amount:  math.NewInt(140600000),
		Claimed: false,
	},
	{
		Address: "elys1zgpjq3lu490jmv0vmp5w6569yc3v6zv8f2ehd7",
		Amount:  math.NewInt(140400000),
		Claimed: false,
	},
	{
		Address: "elys15kwewfemysnkjjqtpcgj2lxj0h0nrwnelc70ec",
		Amount:  math.NewInt(140300000),
		Claimed: false,
	},
	{
		Address: "elys15esydlgz05y94er4jhsnare5mqyz77e5nn57wr",
		Amount:  math.NewInt(140200000),
		Claimed: false,
	},
	{
		Address: "elys1eakhpgn3m5h5k9qwp20pl496h7d4u7f83vuy2p",
		Amount:  math.NewInt(140200000),
		Claimed: false,
	},
	{
		Address: "elys1m8du4mehzfl0y8ps72a6fzazjgg70af0l39f48",
		Amount:  math.NewInt(140200000),
		Claimed: false,
	},
	{
		Address: "elys15uex0vvxjq8nyp3h9ujckwa0jzcw6f3m5sjmnr",
		Amount:  math.NewInt(140000000),
		Claimed: false,
	},
	{
		Address: "elys1tfcy7wvxqar0n4zq26uts42d4vh3vpjer70t2s",
		Amount:  math.NewInt(140000000),
		Claimed: false,
	},
	{
		Address: "elys14v8x0pl7g8qcgc6c08mehehjuvyvt9qsk9vufc",
		Amount:  math.NewInt(139700000),
		Claimed: false,
	},
	{
		Address: "elys1kkyhp3z9c0yfh26aphtxccsl6andlef5vakzst",
		Amount:  math.NewInt(139400000),
		Claimed: false,
	},
	{
		Address: "elys1k7gv40r6qys2cazgm5vjqgu87lx977lqrnmstg",
		Amount:  math.NewInt(139300000),
		Claimed: false,
	},
	{
		Address: "elys10ezdjttdsdr4wcwrmt5yxql57zvkxtys8aynsf",
		Amount:  math.NewInt(139200000),
		Claimed: false,
	},
	{
		Address: "elys12s7cfwfj3jc8f89cw8ws0jfaux37tyjwxz5jed",
		Amount:  math.NewInt(139100000),
		Claimed: false,
	},
	{
		Address: "elys13p25jzyme4sg89e24xs5w5nn52c0tvhjcnnn4e",
		Amount:  math.NewInt(139100000),
		Claimed: false,
	},
	{
		Address: "elys1wyl8y5ewzzyaz4vw4h3z6k8jcsrcmv9hfe3ww0",
		Amount:  math.NewInt(139000000),
		Claimed: false,
	},
	{
		Address: "elys18qz8ecysetqqhyv4qs9tlw5ylpgnx9ymfc5lem",
		Amount:  math.NewInt(138900000),
		Claimed: false,
	},
	{
		Address: "elys1tppd4j6l6fwtm5nyhqhnkpzuf94zv9rpc32djf",
		Amount:  math.NewInt(138900000),
		Claimed: false,
	},
	{
		Address: "elys197mgn28f4qg7h4qk6f6vqtqt2x45j0ap7yvsal",
		Amount:  math.NewInt(138800000),
		Claimed: false,
	},
	{
		Address: "elys1nufpgql23svurhwy4qzrrjze0d8qvt4sw8jl4s",
		Amount:  math.NewInt(138700000),
		Claimed: false,
	},
	{
		Address: "elys1rt3azk2vw9ancze0nxh5w70aa59suel8c4k9v6",
		Amount:  math.NewInt(138500000),
		Claimed: false,
	},
	{
		Address: "elys1lxhevlgwd5n72j8hhkr69z4m7fvhsny833gnra",
		Amount:  math.NewInt(138400000),
		Claimed: false,
	},
	{
		Address: "elys1t94jmzcr9a48g3h48s9edvwf9sp4t6p55h3wny",
		Amount:  math.NewInt(138400000),
		Claimed: false,
	},
	{
		Address: "elys16ck8y8naxk922wj5kv7tuewsj34khe82etp0ly",
		Amount:  math.NewInt(138100000),
		Claimed: false,
	},
	{
		Address: "elys12fn2lcgz3xdxj3fhunpz367aqv6fz6dtlpj899",
		Amount:  math.NewInt(137500000),
		Claimed: false,
	},
	{
		Address: "elys1htlcpzppxhnalgvyljvyd9f67zemwmctspc04q",
		Amount:  math.NewInt(137500000),
		Claimed: false,
	},
	{
		Address: "elys13xm825u8ftc66vahjjctmuqn62m9ap4muvgjj5",
		Amount:  math.NewInt(137400000),
		Claimed: false,
	},
	{
		Address: "elys15dwqhjf3f0shcgtymvpynrahqcqayucrellptn",
		Amount:  math.NewInt(137400000),
		Claimed: false,
	},
	{
		Address: "elys1awjhpe0vuzjzzary0x55dvzwgeqxl5fytdpte6",
		Amount:  math.NewInt(137400000),
		Claimed: false,
	},
	{
		Address: "elys1ljc524mwshz3uuzqgwyvszuxvedm8j6yqfe0tg",
		Amount:  math.NewInt(137400000),
		Claimed: false,
	},
	{
		Address: "elys1ujv2gmfwrwzj504ntggqld0q5euafp76h7g3sf",
		Amount:  math.NewInt(137400000),
		Claimed: false,
	},
	{
		Address: "elys1y9hgg8xwgnkxw5jtelm2klhalvfjep2l52prv2",
		Amount:  math.NewInt(137400000),
		Claimed: false,
	},
	{
		Address: "elys1znq8udruaq2zutmadl5ucyrytk6q44czh8n5yg",
		Amount:  math.NewInt(137400000),
		Claimed: false,
	},
	{
		Address: "elys1aal3unq7g3edp8wth3p8zwdvyn8ygvnq4q05h4",
		Amount:  math.NewInt(136900000),
		Claimed: false,
	},
	{
		Address: "elys1rm3pcrqaf297c9fqep5a92eudchvqmzjlfe0md",
		Amount:  math.NewInt(136900000),
		Claimed: false,
	},
	{
		Address: "elys1eagppq63gml0a9508jqamsn8f0lvxzn0jjwut2",
		Amount:  math.NewInt(136600000),
		Claimed: false,
	},
	{
		Address: "elys1jxfa0vza6h77tclenfsswneh5dklrtkt53whaz",
		Amount:  math.NewInt(136400000),
		Claimed: false,
	},
	{
		Address: "elys1t5pwzc6rltrfsy4m8hnf9r0vx6halp7zsgx43j",
		Amount:  math.NewInt(136300000),
		Claimed: false,
	},
	{
		Address: "elys14sd8h6h7lt2j9ema55036wltjr0gmrm3chzt6r",
		Amount:  math.NewInt(136200000),
		Claimed: false,
	},
	{
		Address: "elys1s9xtexv9wxmfqle304ph9xa3shkw5tsk0tg8yw",
		Amount:  math.NewInt(136200000),
		Claimed: false,
	},
	{
		Address: "elys1pq3ytqwv2syufkkskpenlg868q6vc6pehek9sl",
		Amount:  math.NewInt(136000000),
		Claimed: false,
	},
	{
		Address: "elys12k5c869zd9qmtawz5rutu9njjazk0f79unj337",
		Amount:  math.NewInt(135900000),
		Claimed: false,
	},
	{
		Address: "elys1mwf6yuxlw4z90rtekedjw572ndtar35hjug7zq",
		Amount:  math.NewInt(135900000),
		Claimed: false,
	},
	{
		Address: "elys1ad6ydly9r9x5r5rxdm4zgdzgqv89fqdmw3383d",
		Amount:  math.NewInt(135700000),
		Claimed: false,
	},
	{
		Address: "elys1rtrlp5tk4ql6zuxpy4dyuvh662ut380vxdrjjk",
		Amount:  math.NewInt(135500000),
		Claimed: false,
	},
	{
		Address: "elys14m87h25dehtee4tq9m746uj4zfc7adgtmd8vfj",
		Amount:  math.NewInt(135300000),
		Claimed: false,
	},
	{
		Address: "elys13xqulnkw4maxf0c53g83nhkl3tyxcw7awltzk5",
		Amount:  math.NewInt(135100000),
		Claimed: false,
	},
	{
		Address: "elys1qfa6fzp5u6aclzn4l2x635txdpfhmn93wec55d",
		Amount:  math.NewInt(134800000),
		Claimed: false,
	},
	{
		Address: "elys1uw7cpur9mccjrky3n0v77nmk8kc83eztkk297d",
		Amount:  math.NewInt(134800000),
		Claimed: false,
	},
	{
		Address: "elys1qqw2el8y34zzx4p9qu4cnsh2xvuecex2hnfg4n",
		Amount:  math.NewInt(134600000),
		Claimed: false,
	},
	{
		Address: "elys1wygla84r6fwxsl7j2vzrm9djj383vzlr5nj8r9",
		Amount:  math.NewInt(134199999),
		Claimed: false,
	},
	{
		Address: "elys15cue7ql2rh6kv8kdm6c7d35qnxpu8nemn075nv",
		Amount:  math.NewInt(134100000),
		Claimed: false,
	},
	{
		Address: "elys1jgnkjne6yf9u92e0al7qff2w2fpajxl6lkl5hv",
		Amount:  math.NewInt(134000000),
		Claimed: false,
	},
	{
		Address: "elys1j9g25qslvjvpqjtu0da0yvmumtaqk85u4dveyf",
		Amount:  math.NewInt(133900000),
		Claimed: false,
	},
	{
		Address: "elys1ceaehhr4yu3c40q7pz7d9urw7xsv0a62mpnsym",
		Amount:  math.NewInt(133800000),
		Claimed: false,
	},
	{
		Address: "elys1ctud00kqxwl8xfufkgsvaf7f48e8kn6ekg5hs0",
		Amount:  math.NewInt(133800000),
		Claimed: false,
	},
	{
		Address: "elys1v43p9zstfr9pcrrc7xktqqv0utquuj3xmdrrrt",
		Amount:  math.NewInt(133400000),
		Claimed: false,
	},
	{
		Address: "elys18er0jw7nssafu8jync7qxkru0pl36wvdvt7mq5",
		Amount:  math.NewInt(133000000),
		Claimed: false,
	},
	{
		Address: "elys1ekxgvuu62zlq28snweyw9m6t4umw3jdasc64kf",
		Amount:  math.NewInt(133000000),
		Claimed: false,
	},
	{
		Address: "elys1ha5dcmxrsjte0adtawf94hr3z22pkqd656fntd",
		Amount:  math.NewInt(132900000),
		Claimed: false,
	},
	{
		Address: "elys1gkwxv9sj6r6awztz4tc7dfudx7lzzugu0ls07w",
		Amount:  math.NewInt(132800000),
		Claimed: false,
	},
	{
		Address: "elys1jnmwm0lfvprffl2ugslk9muvkdsq6m5udydmns",
		Amount:  math.NewInt(132800000),
		Claimed: false,
	},
	{
		Address: "elys1qxt4ym2nymln0xhdn420pewqjnalr3ejtfdn74",
		Amount:  math.NewInt(132600000),
		Claimed: false,
	},
	{
		Address: "elys12fdjl5n25avc75g8yh5g93wf9l436y9cnx2rxu",
		Amount:  math.NewInt(132500000),
		Claimed: false,
	},
	{
		Address: "elys16vag927x8g2zlvflqpg5apx5pl56car4kwyfuh",
		Amount:  math.NewInt(132400000),
		Claimed: false,
	},
	{
		Address: "elys1xl3x52f74rcv5ar8lew8ysxdrzy8p7ux7jccfp",
		Amount:  math.NewInt(132100000),
		Claimed: false,
	},
	{
		Address: "elys1eejcn89mtnvzjyn8v74kr9lcye33dmc0z4zm43",
		Amount:  math.NewInt(131900000),
		Claimed: false,
	},
	{
		Address: "elys1g33rf6ugqtrtca68p5dsd2mfh7tsdfu4v6mrm9",
		Amount:  math.NewInt(131900000),
		Claimed: false,
	},
	{
		Address: "elys1gq3rvcngcp3m8lgplahxsx6yw9r2qmj2u0esk3",
		Amount:  math.NewInt(131900000),
		Claimed: false,
	},
	{
		Address: "elys1u3m8fvs5k6rqrazhf6jkmwvlam2vqe8c2s7srg",
		Amount:  math.NewInt(131900000),
		Claimed: false,
	},
	{
		Address: "elys1gt6q2xs2q6e2p6d9aalx72nxfcnl6kyr6eaunz",
		Amount:  math.NewInt(131500000),
		Claimed: false,
	},
	{
		Address: "elys1xp9vpjdzx08uh7llkl46upc3cefmn4kpaj53cy",
		Amount:  math.NewInt(131100000),
		Claimed: false,
	},
	{
		Address: "elys1ggmar2u93aexla65qlqdjthalcvkrscc8cx5tv",
		Amount:  math.NewInt(130900000),
		Claimed: false,
	},
	{
		Address: "elys1kmauuq0weyh45mycwkyd00pmrswytdl29u5t5p",
		Amount:  math.NewInt(130500000),
		Claimed: false,
	},
	{
		Address: "elys1x6huu9wa8yza8ffnes86qx37hglmjmj74974q0",
		Amount:  math.NewInt(130199999),
		Claimed: false,
	},
	{
		Address: "elys1k0a75sar0a5um4fkvpceeap9kmlk9t460zxp3q",
		Amount:  math.NewInt(130100000),
		Claimed: false,
	},
	{
		Address: "elys106gvnnet704z6redzwfd2ufq0mqzaj63h0ql96",
		Amount:  math.NewInt(130000000),
		Claimed: false,
	},
	{
		Address: "elys1z8pryu564uqkplep7azftsjtt5mnf9htja9gyc",
		Amount:  math.NewInt(129800000),
		Claimed: false,
	},
	{
		Address: "elys1mzqww7evte3yh4nduepnq7javg0v9e3tj9p3xr",
		Amount:  math.NewInt(129699999),
		Claimed: false,
	},
	{
		Address: "elys1cquf2vuey34rknph5jt8t34lfw9cj22rct85lf",
		Amount:  math.NewInt(129400000),
		Claimed: false,
	},
	{
		Address: "elys1fkarkse9gkgs20va0dqlcy6lus49kmx4a5hew8",
		Amount:  math.NewInt(129199999),
		Claimed: false,
	},
	{
		Address: "elys1uaucev4kyydwa3x74u5ujq4tq5cemakaparemg",
		Amount:  math.NewInt(129100000),
		Claimed: false,
	},
	{
		Address: "elys1ewqmqzqlrzlu2sv8ryru3znarp6y2vty3qw4s4",
		Amount:  math.NewInt(129000000),
		Claimed: false,
	},
	{
		Address: "elys1p3ye7u6hpx74rcvm5jutl6q7s9uwc0htnz4xj7",
		Amount:  math.NewInt(128900000),
		Claimed: false,
	},
	{
		Address: "elys1nnlvg9fcs7lfsfxgeaf69hg604nxag0axa5ntq",
		Amount:  math.NewInt(128100000),
		Claimed: false,
	},
	{
		Address: "elys1qxgqr2fwzc0pd7mz6vpudlxckffn75ksuspvcj",
		Amount:  math.NewInt(128100000),
		Claimed: false,
	},
	{
		Address: "elys1sk65x8c6s37agpnzsnjagp35zllwnwa4de6xes",
		Amount:  math.NewInt(127900000),
		Claimed: false,
	},
	{
		Address: "elys1y9ya09up4gasjywyh5kkywyhc3dpepxyqtqukm",
		Amount:  math.NewInt(127800000),
		Claimed: false,
	},
	{
		Address: "elys1j28vu2lj8npk04gyuhpw5t70qqehmy3x6emywf",
		Amount:  math.NewInt(127600000),
		Claimed: false,
	},
	{
		Address: "elys1awnu7f9q30mfjevvpxlkx5vp80v4gw5k4ykula",
		Amount:  math.NewInt(127300000),
		Claimed: false,
	},
	{
		Address: "elys108j5umf73a6g8wwgdc6x7d0ups0slyaus24v4c",
		Amount:  math.NewInt(127200000),
		Claimed: false,
	},
	{
		Address: "elys14hj93u2yp2pzm7llx58u6nfe5nf6e03dh3apcm",
		Amount:  math.NewInt(127200000),
		Claimed: false,
	},
	{
		Address: "elys1fqmhzpapp8yr33ezr2txxys2pgvcqmlwxy40je",
		Amount:  math.NewInt(127100000),
		Claimed: false,
	},
	{
		Address: "elys1v45dlhzj7rpgtplr04ytgwksnjxxsnn0egu4su",
		Amount:  math.NewInt(127100000),
		Claimed: false,
	},
	{
		Address: "elys1ck2zy7hf0fyuzjv2f0eujpnr0qz80xwjwssx2d",
		Amount:  math.NewInt(127000000),
		Claimed: false,
	},
	{
		Address: "elys1wlureh2j6exueae8hlp5v67uedjvesjrdr9g46",
		Amount:  math.NewInt(126400000),
		Claimed: false,
	},
	{
		Address: "elys1xqc64za26mrmcfhtsts9hdsu3schpmwaz4cfxf",
		Amount:  math.NewInt(126400000),
		Claimed: false,
	},
	{
		Address: "elys1wmz7zx279uuh88rqvy5f4tfdejg3lq6ks4y93a",
		Amount:  math.NewInt(126100000),
		Claimed: false,
	},
	{
		Address: "elys12d8d2w2q4y2k9jks57rz6th02jhd352uj8w5x0",
		Amount:  math.NewInt(125500000),
		Claimed: false,
	},
	{
		Address: "elys10ja5586gkr8vsml6uetgy8vvwgnke4tx3ha2dy",
		Amount:  math.NewInt(125400000),
		Claimed: false,
	},
	{
		Address: "elys1tr3jssvjyken3cxwrzd3jl2uqd9ne8se4pe5gq",
		Amount:  math.NewInt(125400000),
		Claimed: false,
	},
	{
		Address: "elys1rqpnjx27982n3a9dju08d7zttfuhrgsc6y99na",
		Amount:  math.NewInt(125100000),
		Claimed: false,
	},
	{
		Address: "elys10cptayrxc7xdn0pq872kmd63y45ms8st9c8z3f",
		Amount:  math.NewInt(124900000),
		Claimed: false,
	},
	{
		Address: "elys182x45hyjjes5cxs540lq98zd3z55tl4nvhnsdk",
		Amount:  math.NewInt(124700000),
		Claimed: false,
	},
	{
		Address: "elys1jfwp4d5mg93sapd67r4274f8y8tg4xn758g2rc",
		Amount:  math.NewInt(124700000),
		Claimed: false,
	},
	{
		Address: "elys1ur7tfac3y6m5v4xlvx66avpyv602na4v5safsr",
		Amount:  math.NewInt(124500000),
		Claimed: false,
	},
	{
		Address: "elys1vph0acs9dg695n4cga7d0h7hvhf5ugp2llva8s",
		Amount:  math.NewInt(124500000),
		Claimed: false,
	},
	{
		Address: "elys1a4rv940trkmfmscyfrd8fxez74qm2akp0wgzp8",
		Amount:  math.NewInt(124400000),
		Claimed: false,
	},
	{
		Address: "elys13vn7ja2rtlyjc3raddecgfauwpsq9tplgp4azn",
		Amount:  math.NewInt(124200000),
		Claimed: false,
	},
	{
		Address: "elys177zkqpxasrcjlfkhxknvmsgdkmm593dsq2d6z2",
		Amount:  math.NewInt(124200000),
		Claimed: false,
	},
	{
		Address: "elys18s9jkjrfsn5w4mcmduffgc9fmlq9n3mgfxwcnd",
		Amount:  math.NewInt(124200000),
		Claimed: false,
	},
	{
		Address: "elys10ze66jyp7m0rlvt8aer4hxhgrcunstyjvdfqqg",
		Amount:  math.NewInt(124100000),
		Claimed: false,
	},
	{
		Address: "elys1lz8cdqn6ffuaj0fr9dkzjv9nefzla6ud0u0jau",
		Amount:  math.NewInt(124100000),
		Claimed: false,
	},
	{
		Address: "elys1gyees9j6llcje80lfqf3m3rx3ltcmlw2yjqxtn",
		Amount:  math.NewInt(123700000),
		Claimed: false,
	},
	{
		Address: "elys12zwaujvq38mfd8s35lgwarhhw545zm4e3fdduh",
		Amount:  math.NewInt(123600000),
		Claimed: false,
	},
	{
		Address: "elys1gh2w3wr6plz5z74jftlq7fqvmmgep2t6k3cnkj",
		Amount:  math.NewInt(123600000),
		Claimed: false,
	},
	{
		Address: "elys1dqyny0nal4pa4dd9uaptrdq98yw4awd9gq930h",
		Amount:  math.NewInt(123100000),
		Claimed: false,
	},
	{
		Address: "elys1lstv7rr6wsnaf6rxfk07rdlzr7l0p52c2ak2c7",
		Amount:  math.NewInt(122600000),
		Claimed: false,
	},
	{
		Address: "elys1wcrmk4az7vxrldlxgk8mffxrh09vfunu2sgvs6",
		Amount:  math.NewInt(122600000),
		Claimed: false,
	},
	{
		Address: "elys1yj2hnlsynzkqmhv4hk993svtck0nvqel5zm2qa",
		Amount:  math.NewInt(122600000),
		Claimed: false,
	},
	{
		Address: "elys1mq5u754w234sezjujrdyzruvmu5rqn7z7eqa8w",
		Amount:  math.NewInt(122500000),
		Claimed: false,
	},
	{
		Address: "elys1qupqp7g2xzualda40pelceszkgc5p38qqprthg",
		Amount:  math.NewInt(122300000),
		Claimed: false,
	},
	{
		Address: "elys1svhcykque5zx6vs27dz7hzng45ltvvqlj50qw2",
		Amount:  math.NewInt(122300000),
		Claimed: false,
	},
	{
		Address: "elys1f00sqla2ms8hvtyvl7k9ngk3fnh3648vhz6lg3",
		Amount:  math.NewInt(122200000),
		Claimed: false,
	},
	{
		Address: "elys1gest3cqved5en6ulyqw5pj7tluty5rl3tqkmzt",
		Amount:  math.NewInt(122200000),
		Claimed: false,
	},
	{
		Address: "elys1l9t8kfpv0v87jdlejupz6lvrr58xd32usppray",
		Amount:  math.NewInt(122200000),
		Claimed: false,
	},
	{
		Address: "elys1e5jmcepmgr6wsmp4n3tlhmdcny2pw49cyh2nke",
		Amount:  math.NewInt(122000000),
		Claimed: false,
	},
	{
		Address: "elys187sflkqjtrrc3wz8ly5j93lt3s08036put6n9j",
		Amount:  math.NewInt(121900000),
		Claimed: false,
	},
	{
		Address: "elys1hecjdhsg3tpact9sjeaakfx68lhv562yruget8",
		Amount:  math.NewInt(121900000),
		Claimed: false,
	},
	{
		Address: "elys1kvwfmmvrnduu69fyl92p6wynqzn6cz73yuhsma",
		Amount:  math.NewInt(121900000),
		Claimed: false,
	},
	{
		Address: "elys17e0fpg5ynh3pd8lzwfsenlulrwwam03eyl0hen",
		Amount:  math.NewInt(121600000),
		Claimed: false,
	},
	{
		Address: "elys1pnwyc44sw7twaw53ss42a7uz2qaq4wfmzz6dkj",
		Amount:  math.NewInt(121600000),
		Claimed: false,
	},
	{
		Address: "elys1re7ppjnpx7etwqyhl0esmdpu7egy6e4r6eaj24",
		Amount:  math.NewInt(121500000),
		Claimed: false,
	},
	{
		Address: "elys1uuknxje8vqaad59hldyeclls9vrr9yftyam95l",
		Amount:  math.NewInt(121500000),
		Claimed: false,
	},
	{
		Address: "elys18mwx5sk6755alc7s6te6460c6x5pwn09gq7ftq",
		Amount:  math.NewInt(121300000),
		Claimed: false,
	},
	{
		Address: "elys1nrqgqsplxj8cf6dkf6svutzuukcvvj4hl5dmsa",
		Amount:  math.NewInt(121300000),
		Claimed: false,
	},
	{
		Address: "elys1ru0avp98ruxq6e3qq4jxtfa3ynkvqjhejswtdh",
		Amount:  math.NewInt(121200000),
		Claimed: false,
	},
	{
		Address: "elys1llfpxuxveswlwekvyu479s2dall6u2gfc8jkpr",
		Amount:  math.NewInt(121100000),
		Claimed: false,
	},
	{
		Address: "elys1u8skydc7zllup0w29fm3vgkg69at99nl3ekser",
		Amount:  math.NewInt(121000000),
		Claimed: false,
	},
	{
		Address: "elys1h4xs6wqdqxvtpzc3e8s532sk3sutfwuk6tpuga",
		Amount:  math.NewInt(120900000),
		Claimed: false,
	},
	{
		Address: "elys1t57xnq02h0me3xw8d7qmv3j73s4mj7a7usafxq",
		Amount:  math.NewInt(120700000),
		Claimed: false,
	},
	{
		Address: "elys1z4wee02960mg0k5h4wd2cvyr2cr5mreurgn8dg",
		Amount:  math.NewInt(120700000),
		Claimed: false,
	},
	{
		Address: "elys1dnadzk7zfhv4prsx6cpp3elmq92ngjaeuxp5dn",
		Amount:  math.NewInt(120400000),
		Claimed: false,
	},
	{
		Address: "elys1vmte3lvyyc7he4lfa8twf4fcf2edz6y2hm44q2",
		Amount:  math.NewInt(120400000),
		Claimed: false,
	},
	{
		Address: "elys1hv4rlv94dwwsj3z7l6muye6ndmup0htakxlm23",
		Amount:  math.NewInt(120300000),
		Claimed: false,
	},
	{
		Address: "elys12l438p55kp4458rdrn8uugy6mzs8lm4mc0qa5z",
		Amount:  math.NewInt(120200000),
		Claimed: false,
	},
	{
		Address: "elys1jhdny8ntkassdlww394x5vyh7hvd9lhgqdp4z6",
		Amount:  math.NewInt(120200000),
		Claimed: false,
	},
	{
		Address: "elys1d26hctx0kr827dhm9xw6fgywefrzre9ps5xhe7",
		Amount:  math.NewInt(119900000),
		Claimed: false,
	},
	{
		Address: "elys1wlzxqrtsd7unhx974sy5cxrafrjfdg3c7j2nzw",
		Amount:  math.NewInt(119600000),
		Claimed: false,
	},
	{
		Address: "elys1m8s0rq4hs5ew9eel7lrsjt2hm4803jy8ukfkk9",
		Amount:  math.NewInt(119200000),
		Claimed: false,
	},
	{
		Address: "elys1n6vy9vnq6u5l00s6nv8xds6526f6kulyvreu49",
		Amount:  math.NewInt(119200000),
		Claimed: false,
	},
	{
		Address: "elys1z2qgrklcfarau634ywf7k7vsgnwu5jcrx94l0y",
		Amount:  math.NewInt(119200000),
		Claimed: false,
	},
	{
		Address: "elys18v2ctegf7xl9hvtsyv6swjmqjt3hnnw7jpnxs9",
		Amount:  math.NewInt(118800000),
		Claimed: false,
	},
	{
		Address: "elys1l84dmssk7yxp55ma0xz5zl95f0e2jl5jrcxdd2",
		Amount:  math.NewInt(118700000),
		Claimed: false,
	},
	{
		Address: "elys1qlsjye8spg6858zzsesw9674lgaff5unu3njnt",
		Amount:  math.NewInt(118700000),
		Claimed: false,
	},
	{
		Address: "elys1zm0s4h7yq6h75tpjg20fykjelemjpqp0662nvj",
		Amount:  math.NewInt(118700000),
		Claimed: false,
	},
	{
		Address: "elys1x83phc79jjt40gxq44sekecf462rgdd5rftwp9",
		Amount:  math.NewInt(118500000),
		Claimed: false,
	},
	{
		Address: "elys1f2n7f9c4wpz2qd4he7nlx23kxnz2g64w8l2ers",
		Amount:  math.NewInt(118200000),
		Claimed: false,
	},
	{
		Address: "elys15d23ekfyn0tyszfaaymvdg6v6k0kqn9v0ejdk5",
		Amount:  math.NewInt(117800000),
		Claimed: false,
	},
	{
		Address: "elys1r2nujddqfxvw3nuvdfgsnskxgrdhpfz75a2mtx",
		Amount:  math.NewInt(117700000),
		Claimed: false,
	},
	{
		Address: "elys12nv0wx69u7qfr7yftmtsavvyazu9ggs6vyyhve",
		Amount:  math.NewInt(117200000),
		Claimed: false,
	},
	{
		Address: "elys1etmgwdvzmmqgtqf3wyhjvvncr957tdy0pa8x0q",
		Amount:  math.NewInt(116900000),
		Claimed: false,
	},
	{
		Address: "elys1328u42mjw8lxg45m09635yu3qaanjzndc7yep5",
		Amount:  math.NewInt(116700000),
		Claimed: false,
	},
	{
		Address: "elys1sjldz4dx3u772zz233d0fgzhc79p5vnyujaeku",
		Amount:  math.NewInt(116400000),
		Claimed: false,
	},
	{
		Address: "elys1x5qzjhs6hw52u52h4w9uk57u0lemggjg0y4qkr",
		Amount:  math.NewInt(116400000),
		Claimed: false,
	},
	{
		Address: "elys1h0r5zscgx3ajrg4mcmhqskd3fq0k2v2tqns7qt",
		Amount:  math.NewInt(116300000),
		Claimed: false,
	},
	{
		Address: "elys1a8z62sax7ptdweccuj4uv2s842sk6t4vzmtcq8",
		Amount:  math.NewInt(116100000),
		Claimed: false,
	},
	{
		Address: "elys199ntpdcqpj6pqaev3fffzpyt5yr6xlts8vd2e3",
		Amount:  math.NewInt(115700000),
		Claimed: false,
	},
	{
		Address: "elys1ce0ecx3vfa5dxhnkll8ld2jm30lp0yxpax3uz5",
		Amount:  math.NewInt(115600000),
		Claimed: false,
	},
	{
		Address: "elys1fdyk8z98g0w66mrnue5d0qgxeg585fxsg8x9v0",
		Amount:  math.NewInt(115600000),
		Claimed: false,
	},
	{
		Address: "elys1dmcztf4ru64kuqm8h5c8rqg7l8q2cddu90t369",
		Amount:  math.NewInt(115500000),
		Claimed: false,
	},
	{
		Address: "elys1xwrfl2ysjw7ls69yhr2uklzs8jf2ll7al6wl5z",
		Amount:  math.NewInt(115400000),
		Claimed: false,
	},
	{
		Address: "elys1czw8kd733420g0wpgt7kptprd2ql9dwzvua5hk",
		Amount:  math.NewInt(115300000),
		Claimed: false,
	},
	{
		Address: "elys1mgg6nv9vcv74w63wa68mf48ujc0u2dsr97tkjn",
		Amount:  math.NewInt(115300000),
		Claimed: false,
	},
	{
		Address: "elys16yczfcuvts9yx4fc0pfcplyvudu7p2jgzt9vwl",
		Amount:  math.NewInt(115200000),
		Claimed: false,
	},
	{
		Address: "elys1ja3aznczrv65k4vnutya34kq4r2apsp4gpceva",
		Amount:  math.NewInt(115000000),
		Claimed: false,
	},
	{
		Address: "elys1q2z3tj5yq9cx3nc8gc5tqjk7rrsk4h7cmt4uae",
		Amount:  math.NewInt(115000000),
		Claimed: false,
	},
	{
		Address: "elys1qm09l5ktw6s6exta8marte5rmuhffyswnv3tq9",
		Amount:  math.NewInt(115000000),
		Claimed: false,
	},
	{
		Address: "elys1uyl48e6snaqztpptt0eyme8z656ku4xtdekpf3",
		Amount:  math.NewInt(114900000),
		Claimed: false,
	},
	{
		Address: "elys1yxgkl9zrx6nwpmtr8de3lmhr6t8fs95ptpt7y2",
		Amount:  math.NewInt(114900000),
		Claimed: false,
	},
	{
		Address: "elys1vnkxg6667hrs64nwv9d5g207fd26dz7p4an4n5",
		Amount:  math.NewInt(114800000),
		Claimed: false,
	},
	{
		Address: "elys10ewmtpfhgde2pvjsxvydsmy6nu7uukr8vrraqy",
		Amount:  math.NewInt(114700000),
		Claimed: false,
	},
	{
		Address: "elys1yulpjtsftn8p4g7hx7lczxytkretegr7k3er0f",
		Amount:  math.NewInt(114700000),
		Claimed: false,
	},
	{
		Address: "elys14k5ltledawmsykpl3uj3weh6lg6x96wwrfy5tp",
		Amount:  math.NewInt(114600000),
		Claimed: false,
	},
	{
		Address: "elys17hwsec335k9r6rrtd2fp7pn7uj82gm0f9pprey",
		Amount:  math.NewInt(114600000),
		Claimed: false,
	},
	{
		Address: "elys1qqcagychwcl4gjkyt60lkphmu5q4flsy4dsldy",
		Amount:  math.NewInt(114600000),
		Claimed: false,
	},
	{
		Address: "elys1t2hyx66cecfnn5m82e30ln55kkufw9u3v3gy86",
		Amount:  math.NewInt(114600000),
		Claimed: false,
	},
	{
		Address: "elys17aak2pu9q0kzq6uxerjut5mkh8gce4ycw3e9jk",
		Amount:  math.NewInt(114500000),
		Claimed: false,
	},
	{
		Address: "elys17k065fcptwrt9vp4tzjcmgv5rg7gaggzcgyfvy",
		Amount:  math.NewInt(114500000),
		Claimed: false,
	},
	{
		Address: "elys122ev56u3x6t6v5z7fp67htyh5kpglsyxje4k0p",
		Amount:  math.NewInt(114400000),
		Claimed: false,
	},
	{
		Address: "elys145vs7fa0lxkr8dps5cfs47pty5j69g7w34kddr",
		Amount:  math.NewInt(114400000),
		Claimed: false,
	},
	{
		Address: "elys1aanwt93r0kk6khlr9mqws73rej8a4vuqvmmxhd",
		Amount:  math.NewInt(114400000),
		Claimed: false,
	},
	{
		Address: "elys1aw8d0jjx4nc2a383jzrq46eh7r5vdj7lkc8798",
		Amount:  math.NewInt(114400000),
		Claimed: false,
	},
	{
		Address: "elys1ej0fr5crdppuxvz6u5ydyyy2rpt8gue2vm7h3z",
		Amount:  math.NewInt(114400000),
		Claimed: false,
	},
	{
		Address: "elys1g2k7qzc595dpzpuzlren7kazvl0f3pyryrqju2",
		Amount:  math.NewInt(114400000),
		Claimed: false,
	},
	{
		Address: "elys1h77lz8gt506pm2v9g9wp60llqt6fzu38vz5aqa",
		Amount:  math.NewInt(114400000),
		Claimed: false,
	},
	{
		Address: "elys1rjxz5jk7ty4fh5789qp4kfzzyuyp9t7sqm34q5",
		Amount:  math.NewInt(114400000),
		Claimed: false,
	},
	{
		Address: "elys1rsz009hwhcq35du5tes7qs2zklnmf2cme5esm6",
		Amount:  math.NewInt(114400000),
		Claimed: false,
	},
	{
		Address: "elys1paqmqej7rme7hsljkv3dj7k902xxjm7c7spm5j",
		Amount:  math.NewInt(114300000),
		Claimed: false,
	},
	{
		Address: "elys1vkj232amfmnwyn7hga34eek7389fmhukelk4rn",
		Amount:  math.NewInt(114100000),
		Claimed: false,
	},
	{
		Address: "elys1p5sp522rv3et4vxy80nf082sjhg0pkqw9n8f9w",
		Amount:  math.NewInt(113800000),
		Claimed: false,
	},
	{
		Address: "elys1hr6xsnd0kgrsvwdf40c5vfp9t5dywewy2ua7vu",
		Amount:  math.NewInt(113700000),
		Claimed: false,
	},
	{
		Address: "elys1pa5zq57mgsl8fy6szc4m6th6h3qdyn2mg7vn0s",
		Amount:  math.NewInt(113700000),
		Claimed: false,
	},
	{
		Address: "elys1250hfkls7pu9nwxulzt563udf7nratsd4auykr",
		Amount:  math.NewInt(113600000),
		Claimed: false,
	},
	{
		Address: "elys1myz5ssk30ss8m2e8adsggz24pha4ptygxvjr26",
		Amount:  math.NewInt(113400000),
		Claimed: false,
	},
	{
		Address: "elys1wvsnryy0mtqd3j4gv0pynknlzyvt8d7dpncec8",
		Amount:  math.NewInt(113300000),
		Claimed: false,
	},
	{
		Address: "elys19w7s2lrrq8ceu30pxruxrvevesmukn0ez9atrz",
		Amount:  math.NewInt(113200000),
		Claimed: false,
	},
	{
		Address: "elys14uc5y6tcppyfkjn7u9z0rjkh8z0l7fndkn0539",
		Amount:  math.NewInt(113000000),
		Claimed: false,
	},
	{
		Address: "elys1mdtpe7a87gruntmj9exvvhjt2g6nkw2t0ydh3q",
		Amount:  math.NewInt(112700000),
		Claimed: false,
	},
	{
		Address: "elys1ut4vgw82w34n22eljm83xt32uuflqe6gyddlnl",
		Amount:  math.NewInt(112700000),
		Claimed: false,
	},
	{
		Address: "elys1g2nmav5rqskpvdln2q86enrndkj3tyznecf09z",
		Amount:  math.NewInt(112400000),
		Claimed: false,
	},
	{
		Address: "elys1llmn7qgrd8tn4jf96zxu3mvn5fmatcvvknlwaq",
		Amount:  math.NewInt(112400000),
		Claimed: false,
	},
	{
		Address: "elys1d66cgf0vn4n2zt2tv37qqlngzvaw354kwtst30",
		Amount:  math.NewInt(112300000),
		Claimed: false,
	},
	{
		Address: "elys1a5cjypa77g95jdmvecxwvjv83v0rpx5ua9uq56",
		Amount:  math.NewInt(112100000),
		Claimed: false,
	},
	{
		Address: "elys1ytcfmu9uqzh7lqz9zmu7cvjgcreuecjn9kwmhy",
		Amount:  math.NewInt(111500000),
		Claimed: false,
	},
	{
		Address: "elys1erucnq2xcwc7zn4ey7nhr959glyhdl3xu76qfu",
		Amount:  math.NewInt(111300000),
		Claimed: false,
	},
	{
		Address: "elys1s79pfyjjved0vn4fnzv80srjrs7tp8raz4wuat",
		Amount:  math.NewInt(111300000),
		Claimed: false,
	},
	{
		Address: "elys1m6263u65x66w275pwetkhh0wzpu2x5v6nwukn4",
		Amount:  math.NewInt(111200000),
		Claimed: false,
	},
	{
		Address: "elys1dly8f5zwaya8q4s7j0dj9367lfdc829a5u7rdj",
		Amount:  math.NewInt(111100000),
		Claimed: false,
	},
	{
		Address: "elys19uw8kssmd6fng44m82l3pz67xmnsy66ems0ck7",
		Amount:  math.NewInt(111000000),
		Claimed: false,
	},
	{
		Address: "elys1fc265lvf9h66hujgdnnzmvp8dhgwhg8zzz9epm",
		Amount:  math.NewInt(111000000),
		Claimed: false,
	},
	{
		Address: "elys132mgdg42y8amgew0mfe3sdklysu0dgrlcg8qtv",
		Amount:  math.NewInt(110900000),
		Claimed: false,
	},
	{
		Address: "elys14wth99d0eskcyc80sjzej39g7pcx520wpy29za",
		Amount:  math.NewInt(110900000),
		Claimed: false,
	},
	{
		Address: "elys14z7gafs5jnn4mgsqn5l2l7lzxxgdy249yy8al4",
		Amount:  math.NewInt(110900000),
		Claimed: false,
	},
	{
		Address: "elys1ctlkg69kg6tkkh83vsap53jd2lpq3cm3u4sre8",
		Amount:  math.NewInt(110900000),
		Claimed: false,
	},
	{
		Address: "elys17z786ecvwv2ll4uu9dg39kl84keuvcyrd768nq",
		Amount:  math.NewInt(110800000),
		Claimed: false,
	},
	{
		Address: "elys15xa9znzeq98m0c386smtct69frssxeqxuf3ecq",
		Amount:  math.NewInt(110500000),
		Claimed: false,
	},
	{
		Address: "elys1vryap7hn0aphzc7mwczst7adn29a4z4242d8lk",
		Amount:  math.NewInt(110500000),
		Claimed: false,
	},
	{
		Address: "elys1789keutdwjxr2avjzcs0g5ga5f3p5yv9esswcg",
		Amount:  math.NewInt(110400000),
		Claimed: false,
	},
	{
		Address: "elys1kx9lkh4qk87zcswnkut4thfx98uq2jcdql2dda",
		Amount:  math.NewInt(110300000),
		Claimed: false,
	},
	{
		Address: "elys10peawzyyvvz5u67rwt5fntfwu0tp5hz08rmk9h",
		Amount:  math.NewInt(110200000),
		Claimed: false,
	},
	{
		Address: "elys1rfmz0yhwvmyucr8yjquhltr9hmeyltrke07l7g",
		Amount:  math.NewInt(110100000),
		Claimed: false,
	},
	{
		Address: "elys1z0amxpw29229ylua8vr8s7mgzw5xru6g8qaxgs",
		Amount:  math.NewInt(110100000),
		Claimed: false,
	},
	{
		Address: "elys1uqavktrl3zkn4t0fe0dmwmampta25un2m88anm",
		Amount:  math.NewInt(109900000),
		Claimed: false,
	},
	{
		Address: "elys1s2uf75g30t5ln96pd2dnnevtjvyk4sxyvys396",
		Amount:  math.NewInt(109800000),
		Claimed: false,
	},
	{
		Address: "elys1fqk7pk59kqy2nutrh6jp0r6dluxdpfdj7wjnxl",
		Amount:  math.NewInt(109700000),
		Claimed: false,
	},
	{
		Address: "elys1jgefy4q9grnqvm5vjrclcfafgzs70rmtywgq4t",
		Amount:  math.NewInt(109600000),
		Claimed: false,
	},
	{
		Address: "elys1547hynp58t5kskdcglfec53rpwalvtsly0yfjs",
		Amount:  math.NewInt(109500000),
		Claimed: false,
	},
	{
		Address: "elys1qn5zfa58ty52sklw3vfsv62mk3g3sk57t03t2f",
		Amount:  math.NewInt(109500000),
		Claimed: false,
	},
	{
		Address: "elys12gnnpkz45wrgexek79amx25f4xkje78v2j2u7h",
		Amount:  math.NewInt(109400000),
		Claimed: false,
	},
	{
		Address: "elys18dg4tkhkgq3tzdt0z8penwss2qw2zk9ahc4xhw",
		Amount:  math.NewInt(109400000),
		Claimed: false,
	},
	{
		Address: "elys19gvc68e0am7ljc4q4ndm2vy3ztzku45mzrh8qc",
		Amount:  math.NewInt(109400000),
		Claimed: false,
	},
	{
		Address: "elys1df3lk6m9pjmufxscractxutx3mmlawu8z4d2lw",
		Amount:  math.NewInt(109000000),
		Claimed: false,
	},
	{
		Address: "elys1kpxj63df5h7wuh3c75nvgrastwy05egjkjthkz",
		Amount:  math.NewInt(109000000),
		Claimed: false,
	},
	{
		Address: "elys133wz5medwjss7dzs0z84lha4jjglem6ll29rd7",
		Amount:  math.NewInt(108600000),
		Claimed: false,
	},
	{
		Address: "elys1n5xrurkxlkljgeyuyzvc4cql6kzvf4hdgp0pny",
		Amount:  math.NewInt(108500000),
		Claimed: false,
	},
	{
		Address: "elys1dsh05hnmushsm5jnya2vfcrkgtmks5h3qk482y",
		Amount:  math.NewInt(108400000),
		Claimed: false,
	},
	{
		Address: "elys1ezkeya075hgeqy8wgul0k3pl5hvukxag08rk8e",
		Amount:  math.NewInt(108400000),
		Claimed: false,
	},
	{
		Address: "elys1xfy96fx88wge8l7keepnqcwqrlpv39ctdrzgky",
		Amount:  math.NewInt(108100000),
		Claimed: false,
	},
	{
		Address: "elys1qqurjlta2xs0py035rcavxy7dtpjg5jhuz9rlu",
		Amount:  math.NewInt(108000000),
		Claimed: false,
	},
	{
		Address: "elys1wjmcgkqwz0gp67433rlhe8ft4cpdhanxkw5rdv",
		Amount:  math.NewInt(108000000),
		Claimed: false,
	},
	{
		Address: "elys1j3j5aff5nxknscqyjc4t6pypjn9tyyeegsqwq9",
		Amount:  math.NewInt(107900000),
		Claimed: false,
	},
	{
		Address: "elys1ux79zz6ywd5u5r29a2y2wf8d6ar8fdcph5un8u",
		Amount:  math.NewInt(107800000),
		Claimed: false,
	},
	{
		Address: "elys193k82lrft8dr95rztfwn44etmzkmzp2vmg8t2q",
		Amount:  math.NewInt(107600000),
		Claimed: false,
	},
	{
		Address: "elys1vt0q5u27j7n39zd6hm30p7csm2gdkmma780m5z",
		Amount:  math.NewInt(107500000),
		Claimed: false,
	},
	{
		Address: "elys1hceylmlseee0v54r5e8lls854rh5qrqex6ydeg",
		Amount:  math.NewInt(107400000),
		Claimed: false,
	},
	{
		Address: "elys1zumu5f9zewmfmqt4gn7neu4cca62xtscmzcvs0",
		Amount:  math.NewInt(107200000),
		Claimed: false,
	},
	{
		Address: "elys1326x2eutldzsrhu048x58zqgtwymg0hy76tlec",
		Amount:  math.NewInt(107100000),
		Claimed: false,
	},
	{
		Address: "elys1c2zqp6289tgcl2zf4lg54dpv7kr0f2xdqvcwp4",
		Amount:  math.NewInt(107000000),
		Claimed: false,
	},
	{
		Address: "elys10t4j8eld3x8mrrkfste53q6k9h3al4stqke0zy",
		Amount:  math.NewInt(106900000),
		Claimed: false,
	},
	{
		Address: "elys1zz7svzr8hvalzd2t9x6a8ufh6u0sjxp329erv0",
		Amount:  math.NewInt(106900000),
		Claimed: false,
	},
	{
		Address: "elys1udagr8pl8wcm45maf7evk8xdcys604j3cjwm9f",
		Amount:  math.NewInt(106800000),
		Claimed: false,
	},
	{
		Address: "elys19jkvphz7yx6w8yy9ksdlm4jvmhrcctw5dts3ah",
		Amount:  math.NewInt(106500000),
		Claimed: false,
	},
	{
		Address: "elys1gf9yy3e3pa6qzhkq09j5g9v927mtpald053v0n",
		Amount:  math.NewInt(106500000),
		Claimed: false,
	},
	{
		Address: "elys1vh0w9edca0kvefsf32vvqqjmjq70m2jhusqj4y",
		Amount:  math.NewInt(106500000),
		Claimed: false,
	},
	{
		Address: "elys1ssksjn3ptfwsmxkv6d2kyjalm478sx5ffe36ln",
		Amount:  math.NewInt(106300000),
		Claimed: false,
	},
	{
		Address: "elys1yd6fdf62z0vk9f75h79qntyvky9hdynn9j8lxz",
		Amount:  math.NewInt(106300000),
		Claimed: false,
	},
	{
		Address: "elys1py4wyc4g3my0v6tl3yyhh28kyfq9mqarn2j3gz",
		Amount:  math.NewInt(106200000),
		Claimed: false,
	},
	{
		Address: "elys1tfkfkhd3lcte0kpttcm84yc952vcshth6rq860",
		Amount:  math.NewInt(106200000),
		Claimed: false,
	},
	{
		Address: "elys1rtgvhfuvhrqasyp9yzs4x7sy8e828r02vze0u6",
		Amount:  math.NewInt(106100000),
		Claimed: false,
	},
	{
		Address: "elys1cw37lg8u5aq2ry2utp4y20vgxc52affwaxudk3",
		Amount:  math.NewInt(106000000),
		Claimed: false,
	},
	{
		Address: "elys1ye05tsdhgn6mdfwetcw4fzwltl6zuw4j697tvn",
		Amount:  math.NewInt(105900000),
		Claimed: false,
	},
	{
		Address: "elys1uwjetn6k0hexwe7thg9pw7tjp3ees9msrsejt6",
		Amount:  math.NewInt(105800000),
		Claimed: false,
	},
	{
		Address: "elys158haaud2zmrm84l2dqfd27q4reskd9cy7ek0ac",
		Amount:  math.NewInt(105300000),
		Claimed: false,
	},
	{
		Address: "elys15qw8uw8weypv923y5kgx43s7ekrue5jctm4urn",
		Amount:  math.NewInt(105100000),
		Claimed: false,
	},
	{
		Address: "elys15y233ekrw7p60kxhzf5pa6agkvt6kk8s0jaz5c",
		Amount:  math.NewInt(105100000),
		Claimed: false,
	},
	{
		Address: "elys1ea8yng52alk8ap9xjfun7jhmpmyl3xgqeztvrj",
		Amount:  math.NewInt(105100000),
		Claimed: false,
	},
	{
		Address: "elys1j94u064l0r60m9wjdj9zkfrchwy434y7h8xn4l",
		Amount:  math.NewInt(105100000),
		Claimed: false,
	},
	{
		Address: "elys1ja26yw0jjdkdrwyjs0qgqkcs9getyj0clgcp8k",
		Amount:  math.NewInt(105100000),
		Claimed: false,
	},
	{
		Address: "elys10datnnlcjmrdl37ka0g4u83chvxpfafmhdmwrt",
		Amount:  math.NewInt(105000000),
		Claimed: false,
	},
	{
		Address: "elys106k968y8xdehfdw6yn7mayr3f03p2tc5ry4wgy",
		Amount:  math.NewInt(104800000),
		Claimed: false,
	},
	{
		Address: "elys14u73spz729z7v8hepa3kp5vd5rkys65jf4p7m2",
		Amount:  math.NewInt(104800000),
		Claimed: false,
	},
	{
		Address: "elys1j0q259tj0mmww8ffe3lvuwzrkufvf9tan99vu6",
		Amount:  math.NewInt(104800000),
		Claimed: false,
	},
	{
		Address: "elys1ru6grcdzh5fckslelv5ysf0l9qaasufryl5g09",
		Amount:  math.NewInt(104800000),
		Claimed: false,
	},
	{
		Address: "elys1xm4vq70p5f2a8e6sfx8rz7j59jd7wcclcqwl3f",
		Amount:  math.NewInt(104800000),
		Claimed: false,
	},
	{
		Address: "elys1vyrctvn55pg5m4e4xsg2dmqh4hwjs3xumylwpj",
		Amount:  math.NewInt(104600000),
		Claimed: false,
	},
	{
		Address: "elys12sj96vljjkyxalqmpyl89qmwn3mclcvaqk8f5r",
		Amount:  math.NewInt(104500000),
		Claimed: false,
	},
	{
		Address: "elys1ng40tj3zlr2dyn50sv3r0m6glx7t39t8v2y3x0",
		Amount:  math.NewInt(104500000),
		Claimed: false,
	},
	{
		Address: "elys1fgtp7k26enhja3w6cq37yxxmsfugczpk4zxmn2",
		Amount:  math.NewInt(104200000),
		Claimed: false,
	},
	{
		Address: "elys1hjdjjt3lh8hqvh3vhmum20evqwm674ey5g6a3j",
		Amount:  math.NewInt(104100000),
		Claimed: false,
	},
	{
		Address: "elys1mdjewj42whquvay093zlfm3lnuh8w9l62ksv42",
		Amount:  math.NewInt(104100000),
		Claimed: false,
	},
	{
		Address: "elys1r7wzfnz2p8f0xwepudg0rr9wpdrxefz6nte2r6",
		Amount:  math.NewInt(104100000),
		Claimed: false,
	},
	{
		Address: "elys1tv8wn039442yqp05r3st739l3uacweatr7xcxx",
		Amount:  math.NewInt(104100000),
		Claimed: false,
	},
	{
		Address: "elys1v2wmsgvu3c7ss3altqqrxhsen8lwenttqkkt9j",
		Amount:  math.NewInt(104100000),
		Claimed: false,
	},
	{
		Address: "elys1m22jczq9r20ahgnhj5f9yqkm0fuhcr80lmzha5",
		Amount:  math.NewInt(104000000),
		Claimed: false,
	},
	{
		Address: "elys1g7xpg3shujrnxqlgr2y5e8kacwayjycm2rdc5w",
		Amount:  math.NewInt(103800000),
		Claimed: false,
	},
	{
		Address: "elys1eq0nwqc2rrkfx3rc8ergvzvae73u93r0r89mpg",
		Amount:  math.NewInt(103700000),
		Claimed: false,
	},
	{
		Address: "elys1l32l47dp5x2zdz9ndcnncz7rwzdr02gnzqnpmg",
		Amount:  math.NewInt(103500000),
		Claimed: false,
	},
	{
		Address: "elys15dkppelz49pryhp4uw6m2hhu7eq5eyuth2j4rk",
		Amount:  math.NewInt(103400000),
		Claimed: false,
	},
	{
		Address: "elys17mrgrpdyt2e5a35zr3ej0474s2j8824ccxyutw",
		Amount:  math.NewInt(103400000),
		Claimed: false,
	},
	{
		Address: "elys1af4pg2mzwxftpmupz8le04c3cw2p3n02zqnxax",
		Amount:  math.NewInt(103300000),
		Claimed: false,
	},
	{
		Address: "elys1rj9xumcl7wedy0fshnazapwke2vtd0dhpww77z",
		Amount:  math.NewInt(103300000),
		Claimed: false,
	},
	{
		Address: "elys1w39lutxau9gn5yn78napfytu3v8d0hg6nhm66e",
		Amount:  math.NewInt(103200000),
		Claimed: false,
	},
	{
		Address: "elys17dtccqyxwsztn96qgxzezatz8n8vctqnmuxh6d",
		Amount:  math.NewInt(103100000),
		Claimed: false,
	},
	{
		Address: "elys1j79dw5yjuqhhqwx29mfzkenz8p34qnt5vs9d6l",
		Amount:  math.NewInt(103100000),
		Claimed: false,
	},
	{
		Address: "elys17ak5yzgmsmexq9ezy84lrcwny3c04g0pdxedaw",
		Amount:  math.NewInt(103000000),
		Claimed: false,
	},
	{
		Address: "elys1e848rnwawj6k8wl7d2vw0u8af8xr7gsf6xnc42",
		Amount:  math.NewInt(103000000),
		Claimed: false,
	},
	{
		Address: "elys1jgal56l73vgw9uzq2au5lm5fp9v6zlyrscqtu6",
		Amount:  math.NewInt(103000000),
		Claimed: false,
	},
	{
		Address: "elys1nn8q5zp4jkvgfy3m7c9jktzxa3hdtudvkgc68n",
		Amount:  math.NewInt(103000000),
		Claimed: false,
	},
	{
		Address: "elys1s67dxfqx2wnhjpd4w27dw9cn734fa4l08h3sg8",
		Amount:  math.NewInt(103000000),
		Claimed: false,
	},
	{
		Address: "elys1xqlc5sdhkhw8p6c3ezm0vxpraavmvg6xdvjyly",
		Amount:  math.NewInt(103000000),
		Claimed: false,
	},
	{
		Address: "elys17u34javanxdndzpd63e5uuudtrclvwsv2k92wr",
		Amount:  math.NewInt(102900000),
		Claimed: false,
	},
	{
		Address: "elys1t9zfkacmlvgz2plm0kau58w7hxklflzavgxwcm",
		Amount:  math.NewInt(102700000),
		Claimed: false,
	},
	{
		Address: "elys12yqx6ej76nsjcrgypmxmzqe905l8n2rx08munt",
		Amount:  math.NewInt(102500000),
		Claimed: false,
	},
	{
		Address: "elys14935y5cyd7vaezqlcn35nhcc9kpmzlvtnpn4p6",
		Amount:  math.NewInt(102300000),
		Claimed: false,
	},
	{
		Address: "elys14wes3e7s2q547t03tdul5h5dq7npeewq9sf3aq",
		Amount:  math.NewInt(102200000),
		Claimed: false,
	},
	{
		Address: "elys19ky9fj9g4wr64n2wwlrwpwyv3nfd9f2ums3qzk",
		Amount:  math.NewInt(102100000),
		Claimed: false,
	},
	{
		Address: "elys1htzeq4wagzcynesp8ejqzhdhv3p7gw62pawwlz",
		Amount:  math.NewInt(102000000),
		Claimed: false,
	},
	{
		Address: "elys1239m55933tzrafgkdq543klffku0us37ywkq49",
		Amount:  math.NewInt(101800000),
		Claimed: false,
	},
	{
		Address: "elys12hhr8heg4vgnsl4kkqns24tx3wspcu3m745tld",
		Amount:  math.NewInt(101800000),
		Claimed: false,
	},
	{
		Address: "elys18t7xyc7yu77a25nnxma059vp7lvk0zt7tzf05t",
		Amount:  math.NewInt(101700000),
		Claimed: false,
	},
	{
		Address: "elys1wtldn72sklaqfy6nygy6tst3taaj6fyz4xwvs7",
		Amount:  math.NewInt(101700000),
		Claimed: false,
	},
	{
		Address: "elys1cnm0dmkpnxlfmajx8542yt0wezrdwj9y79p77z",
		Amount:  math.NewInt(101600000),
		Claimed: false,
	},
	{
		Address: "elys1jzq49x4uy670xr8q84lthaelm62f973qpxjkay",
		Amount:  math.NewInt(101400000),
		Claimed: false,
	},
	{
		Address: "elys1j3hwqxn0t7mx08l2a2wdw7xqtswp0ywf062aav",
		Amount:  math.NewInt(101200000),
		Claimed: false,
	},
	{
		Address: "elys1jf0zhtqxw2922e7caz47jx8tejltdr5tk62s54",
		Amount:  math.NewInt(101200000),
		Claimed: false,
	},
	{
		Address: "elys1syv0r38yy662uz6yqhwpywm6hg0vzuerem93ax",
		Amount:  math.NewInt(101200000),
		Claimed: false,
	},
	{
		Address: "elys13p8cj5ky57rglwwtnh073sfjkxddpen6plnthu",
		Amount:  math.NewInt(101000000),
		Claimed: false,
	},
	{
		Address: "elys17ypm0zx6tfm5zf49kpj3dpyvm6xz9vff398c9u",
		Amount:  math.NewInt(100900000),
		Claimed: false,
	},
	{
		Address: "elys1yn4d53apq5jsfcqpul5w36l332p6yg6x7q9w6c",
		Amount:  math.NewInt(100700000),
		Claimed: false,
	},
	{
		Address: "elys1z3u86ecz0ayqxghcnsdsn5zqjnmkx5upjdxjyh",
		Amount:  math.NewInt(100600000),
		Claimed: false,
	},
	{
		Address: "elys150sfwetjgnaa3nxa0xytv9f8ctlc0hu8yf7yaa",
		Amount:  math.NewInt(100400000),
		Claimed: false,
	},
	{
		Address: "elys1cjd8stpxeu42tqpk50dkl2ufvm20qxcp9qkq2s",
		Amount:  math.NewInt(100400000),
		Claimed: false,
	},
	{
		Address: "elys1l6u22tnv0h9lklc5789cxnjqwuhq7fn9je5h4h",
		Amount:  math.NewInt(100300000),
		Claimed: false,
	},
	{
		Address: "elys1n42lyhmtlzn38zwsu5shcwt5k3ncn4a9sz2zgl",
		Amount:  math.NewInt(100300000),
		Claimed: false,
	},
	{
		Address: "elys1knufu7lftm5pu28up08amy5uhpj5r0ma4g8kah",
		Amount:  math.NewInt(100200000),
		Claimed: false,
	},
	{
		Address: "elys1urgxmkx2wra6adrx7chlql342dgqrqatunp0pq",
		Amount:  math.NewInt(100200000),
		Claimed: false,
	},
	{
		Address: "elys1vc0vwdlyfuuz9aduqehjyc4mzgd57gud7ejgf9",
		Amount:  math.NewInt(99900000),
		Claimed: false,
	},
	{
		Address: "elys1cax835d2jf8422njjynya6fyx3yxyng47lw6ty",
		Amount:  math.NewInt(99800000),
		Claimed: false,
	},
	{
		Address: "elys1q8jyeqygw5cdslgqnqqdszd0xwmdxehqj8l23v",
		Amount:  math.NewInt(99700000),
		Claimed: false,
	},
	{
		Address: "elys1svaj6q52n77rdkl8x2vdudwj846c7ntvx685xx",
		Amount:  math.NewInt(99600000),
		Claimed: false,
	},
	{
		Address: "elys1z3hhg3ngrl6ld7g5gcr8q3na8lpe5ly8xe3csd",
		Amount:  math.NewInt(99600000),
		Claimed: false,
	},
	{
		Address: "elys1a9arttuxfffwyf3cl9jlfp83tfxk2z24wxyeh3",
		Amount:  math.NewInt(99500000),
		Claimed: false,
	},
	{
		Address: "elys1v4nt62sqsplaqcx75w7t96z94lhez6xsl35809",
		Amount:  math.NewInt(99400000),
		Claimed: false,
	},
	{
		Address: "elys18mhf35pmch5ulpwnzx0uhvxzachma2u0sx5gas",
		Amount:  math.NewInt(99200000),
		Claimed: false,
	},
	{
		Address: "elys1m8rwxusczmlv9j5aw92dydf8lvcare3kjapmrg",
		Amount:  math.NewInt(99200000),
		Claimed: false,
	},
	{
		Address: "elys1e55egljfk6vlj8gnashxxszqmmjh35x62uza9j",
		Amount:  math.NewInt(99000000),
		Claimed: false,
	},
	{
		Address: "elys1l4flznyd2kxx8llk0c7slga2malmxsd78c4sjy",
		Amount:  math.NewInt(99000000),
		Claimed: false,
	},
	{
		Address: "elys1vd8gh5ahwnyne0qaxnnstme6nk9ae854kscdjx",
		Amount:  math.NewInt(99000000),
		Claimed: false,
	},
	{
		Address: "elys1kmyct5ajjp7v2pdhjcakd26u4p6ek5uk7s7je0",
		Amount:  math.NewInt(98900000),
		Claimed: false,
	},
	{
		Address: "elys14wxwl3yr0gf496xysy20nkz3g9asz5nzkp9rka",
		Amount:  math.NewInt(98800000),
		Claimed: false,
	},
	{
		Address: "elys1uqmtf087gupq28gtdsh8a53hj3ph33pee768wq",
		Amount:  math.NewInt(98700000),
		Claimed: false,
	},
	{
		Address: "elys18rn934dfk29y2a0csmrdrkakpclls9rgugfte2",
		Amount:  math.NewInt(98400000),
		Claimed: false,
	},
	{
		Address: "elys1l40spywd80kglxs7a2axza7lrzt00pe84x6a6r",
		Amount:  math.NewInt(98400000),
		Claimed: false,
	},
	{
		Address: "elys1cd62pp6vdjwv3fqumwfsuas5xwqvudcr7q7yjy",
		Amount:  math.NewInt(98200000),
		Claimed: false,
	},
	{
		Address: "elys1vncfxpjarteglylmawp7x2qnc5kltw0yrjk9mw",
		Amount:  math.NewInt(98200000),
		Claimed: false,
	},
	{
		Address: "elys190a3de8tgga637f4wgmkwdyfw2dp8zae6nj7k0",
		Amount:  math.NewInt(98000000),
		Claimed: false,
	},
	{
		Address: "elys19daxad8zlclv8nlmnkyfumj76hl7uc7svr80gz",
		Amount:  math.NewInt(98000000),
		Claimed: false,
	},
	{
		Address: "elys17nun6cg7a5xtu9fz2p8scxr7u9jmrf5vv56e6w",
		Amount:  math.NewInt(97700000),
		Claimed: false,
	},
	{
		Address: "elys1xllgfgjahhvr60l8e65y2j9scjldn8d6tjf48p",
		Amount:  math.NewInt(97700000),
		Claimed: false,
	},
	{
		Address: "elys1zv8dwrsv7nkgs2tm3aknssgl4tex8qcy3prdp5",
		Amount:  math.NewInt(97700000),
		Claimed: false,
	},
	{
		Address: "elys1h0ytgclzxf8d6pm45hf8672gs3hy8ptyd6dut4",
		Amount:  math.NewInt(97600000),
		Claimed: false,
	},
	{
		Address: "elys1l8ja40s7vrzfe7k4yucfn3kd6tdau7vud0nj48",
		Amount:  math.NewInt(97600000),
		Claimed: false,
	},
	{
		Address: "elys10txn8xjec6q6rd2ldykehgytm42eap47x6hxv2",
		Amount:  math.NewInt(97300000),
		Claimed: false,
	},
	{
		Address: "elys1za9j2qmu5enw2xcajkqdgetpg84ma0ljgecwnp",
		Amount:  math.NewInt(97200000),
		Claimed: false,
	},
	{
		Address: "elys1jetmwxz5duqghe8gd7t4vk2nluc7zm0afgj4f0",
		Amount:  math.NewInt(96900000),
		Claimed: false,
	},
	{
		Address: "elys12uuweyuphf6mlajazkpxke2hxwdfew2wqaga4q",
		Amount:  math.NewInt(96500000),
		Claimed: false,
	},
	{
		Address: "elys1066cux6dt7ddujwjw05kkzncf6w702a6w52nnk",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys12dffduqxup4n3wd5ag4dkj2exnj5zpcwrz3999",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys12w9gse7lgptef2s336ksjwkuczh3w8nt2ptl4g",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys13n4nzn32gjlhtlmn74yh4jthafjpz2cgzrwfgs",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys13na66c073zz37n83lclwx3tks546kceq3s55c5",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys13qmf79yujwe7t3fllv03sx3c9jmnvgmnu8ea6m",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys17tzxd64nesfueye4dwx7htu5vk2c9hjcl4nkwx",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys19za86lgeqlaqv784pgyhqvuwy0hs8rq2xdns83",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1c2e9xzs7ecwymx9ywnpu5murx6wgvjtvjzxwmj",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1cf4ws4thsrnx3f02da8nzu04l24mcenf3ceawc",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1ddlaplql7ykzmwj496fm9gdlkrd2j89jw5p8av",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1dw9yqppsnrfa6nkgsa7ls0meyfg9qyc5vzwcf2",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1ezdpqwa8uyg78y3cn9xpmn9yaqzudjyvqahh3g",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1ezjqrepzvps402rrt2qgek3yelpqdcx9yrndgg",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1fvqz4g7dg5xf4sa3l75sftphcuwd8lnv65dqwg",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1g26qk8vqrarlv0jjvy70tss4su83n0djpyhpua",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1g5kv40t2hnr3ltj7n00vvmu37cm5u6ekxngmhy",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1g78x3ttr3r4f94sanrnvzt4zzw4vypf72gg37j",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1hatw4r0jfn5pnuedmvyrpcksjfnr7p3n2800jy",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1hhmq6ztrm8xy45czykaje8jktrt63y8cmx4x06",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1hj8n5w8khc8q78f6th627mj6xp25mfzy3ljgj5",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1jn40y9sfe8qa34wzzn7lcpmlfaylvf6pvxwqsm",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1kaqnvnwzt9tcku89v5rtjvvlpny77v5q36zsgw",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1kmc0m6s2aprw8wwlrktypc7mmjcdlnyak6ct99",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1mvlsxwaamefju4tlv7qwfrklecllf0q537n6yf",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1q6xr9zv6svnaj3978cl46055lv4l3h9fmj4j9p",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1qvncs5hqcf9zctkjzsgr2kmw3wkzf69kdwsqy2",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1r9u46y6qmwl4nuvpfptgmg6zxx84akxn96slh6",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1tkvpven4e7ch8hlvtk8ttjk8dvxhyztj305ty4",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1twlcf9e2wlwjh5e3cnv0cckjje02u3adwxs8ur",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1uy2xuwtf7edvek22vul9vqhjlj7cpwxgssaweq",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1vhl5jhs8tx6ret0ynjkjnejp9sja5gjnqc28es",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1wexlwn345qvklwk92lhqx6709pm56gwdhwrwjq",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1wt0zd0wflhccv5val5y97fq4g9c0p2t4tvvz78",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1yd47vuuc7ra0d6sl34d5pzgkvx7dsr6smz44f0",
		Amount:  math.NewInt(96400000),
		Claimed: false,
	},
	{
		Address: "elys1frayjk22hgtszxxk8mxmrsf2eahdflz52suu5k",
		Amount:  math.NewInt(96300000),
		Claimed: false,
	},
	{
		Address: "elys1wuwl5z7k55mmgg507vr9sm2ee84at8hcz69xue",
		Amount:  math.NewInt(96300000),
		Claimed: false,
	},
	{
		Address: "elys1xzz0ukyv0rkwfvxvqp3wd5rkym6p8etk7nev4x",
		Amount:  math.NewInt(96300000),
		Claimed: false,
	},
	{
		Address: "elys1zkxn96zkgv48ptuhjjzpgdjwdd6tqrpzjdn4hk",
		Amount:  math.NewInt(96300000),
		Claimed: false,
	},
	{
		Address: "elys10m6ll6lkdrl7hxgfgjzy80t9f77d4056xyzzum",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys12trcwfd0e82g6904ntgka7y72efjp82hwr70tc",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys12x920cjr7dna5962kstplye2lu568ga6n42fsd",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys12xz4avr9vax6k9p5cufg4j7rtela960jkt5ngx",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys13ltnycyyxgyckprj06qsgeal4mz85k0znfpx64",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys13m580ykcc5yw885y0wfkqjf064z587huzpquk3",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys15hcz0k4lctgvnmddfl4489fyrt9vkdalre5aju",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys169try2l07x6j2slkvhku65h6wdrestw6j80998",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys16sf7swmwum0lgj042e82xzfu4jwrjg5lvll0h5",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys17wurujn0nx6fl0639zjldwakxvmycscsfk3jhx",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys180czrhezntj9cnlf2dg0zk5nr28jgpxvv9n474",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys180gce0efv2nxh3srwm9vmcm5lfv4le86vxj64v",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys18hqlg4wrzulyacw4zmggaufysawcwksdf9fywf",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys19ka0gq4p6x38rsvnqlv0z202udzykdu3ls2pyy",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys19p34gnt0l9s66apjj49zckqefw47pe5x0we5dl",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys19zvd24xd75majwvzn393jwhuctdum47gfq4chr",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1aq4kaqmpwyr2upxuc0h302h4hy5q85p2qthse2",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1c4n7093uqxje2alpa3jhgq77lv7mfg8lnndl9s",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1c5xxqgncr9ymkqfsfsv6aklta636ytyazrp6av",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1c9y65rdpxruxjl3lhf6r0et67r064hymytemhq",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1ckmu4zhxfqv9q5mkfddwxg887tpya4hpfuzh62",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1cn7mjrwhspz3qxqff9l4p3h6qpspvsefuayzhg",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1ctah2fx2r645pllq79k2ftt88ums4nx6nxghnh",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1defy8ws73hhpha7k0hqyumz3gm6fx0euwdlw9h",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1dgw6pfqlvn2d4fhc9yvng3uj54jhxwqefnnq0c",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1fmufm0p2ayhpyp4wmhy4tfxc3elay0tzpn4yth",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1gsx24x8yhx3pm2cl6h9y685kvsqa6fs84fn7hu",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1j6g9najmvgwtjjjt58ddrew8ge4jfhc7wnagm9",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1ku7zgpjerqqfreg7qp4ds75gw86ghcp4waxc95",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1mram5gaa8q5p6m7rtxqu9resh9qckrlp548wxz",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1n2usnhajhy0sdcf3660rss86agptxuvplsdrnu",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1n2uz85m8y847tj2j066jq9kpjvqufj7zn09cz4",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1ngau0rwzh3enl24g7gfpp6kmjyc2wmeyx3khk7",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1p070dt2lpnye6aprmugll0cm8vjx54epfy04rk",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1p8h0vdeelhkw3r5rhagdt9968przfh3dwl52xf",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1pe4a7ucjawyd6dqlnudmv270zpetu2rgmyxfts",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1q4qkla60kmjv9qjpc8avgvzwlqg56x65e8vz8x",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1senj32slk5ne932jl35ek537pwn3utlsva06wh",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1sz8xpdr2jj4q42dk8jk48yksg78aa4vyjlu3n7",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1tj7ry70376dj6m6r40l4dpupk5wlnl83vuk69r",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1tr55zghf7d6rh3l4qfjhs4p4la2pz79z8dsgel",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1uvuu9lzayqf5j54cu7tnx5q0gv7x3le9chzyge",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1v7j3vgv69rzg26laqfqp64us7mwxjg2svwahm4",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1vted30y8u47h03v465e55cv3r5k2tvvzkjfzhx",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1w7dk46xl80krav8dq50q6tap68uvxevgm4qrvm",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1xk9g4ds39llkl95lwywtm79xqf2drpky9gw5qk",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1zns34z73gnu58lqxpmgadt8h4ws9fa4j3ny8r6",
		Amount:  math.NewInt(96200000),
		Claimed: false,
	},
	{
		Address: "elys1063rhg4ske6q8urhk82gpw69ta62ptwldahqmn",
		Amount:  math.NewInt(96100000),
		Claimed: false,
	},
	{
		Address: "elys1flr8qwd897wph2s8tx7l376a7aa2rvh5wylt65",
		Amount:  math.NewInt(96100000),
		Claimed: false,
	},
	{
		Address: "elys1p05mqpfv8r3mersr3xez6eaaerw8h6zuhdgn23",
		Amount:  math.NewInt(96100000),
		Claimed: false,
	},
	{
		Address: "elys1vzvmfq848r9w9nhjx5sa3pmudhg7vtftku2k75",
		Amount:  math.NewInt(96100000),
		Claimed: false,
	},
	{
		Address: "elys18es0wghvc9dpv8f787rpwjzwkn2m2xfx5nf0zc",
		Amount:  math.NewInt(96000000),
		Claimed: false,
	},
	{
		Address: "elys1cvvvh6dudknpcvct9c30r8ss2tzdj0yn22tp03",
		Amount:  math.NewInt(96000000),
		Claimed: false,
	},
	{
		Address: "elys1k3r046f2vaelzyp2wywuvy6ngt68utwl3pyzwl",
		Amount:  math.NewInt(96000000),
		Claimed: false,
	},
	{
		Address: "elys1v3hjp6kyyx28uprrqu3kus7cg3cywj2gspedyj",
		Amount:  math.NewInt(95900000),
		Claimed: false,
	},
	{
		Address: "elys139dudhyk4lqj38axpu23phcfklgkvamcepq20j",
		Amount:  math.NewInt(95800000),
		Claimed: false,
	},
	{
		Address: "elys1dd70xegja3cst7ptt9893wa94kaxfj0lpqydlf",
		Amount:  math.NewInt(95800000),
		Claimed: false,
	},
	{
		Address: "elys1ff93qj0ka5hatpvyk4zsnn8qxfek632z80qnnr",
		Amount:  math.NewInt(95800000),
		Claimed: false,
	},
	{
		Address: "elys1fxc9rrf7lrdy909ch0r8a8zhtasas7p0sv327z",
		Amount:  math.NewInt(95800000),
		Claimed: false,
	},
	{
		Address: "elys1kh0thnn6t04xk55ya27l2gufpdevy0vlzy9yxy",
		Amount:  math.NewInt(95800000),
		Claimed: false,
	},
	{
		Address: "elys1l69s59zhzn0m0035wa07tek5crqmtm27zqt59x",
		Amount:  math.NewInt(95800000),
		Claimed: false,
	},
	{
		Address: "elys1lsw93n5xt67q3nug6ldutcw50vtm9spxaf3kx6",
		Amount:  math.NewInt(95800000),
		Claimed: false,
	},
	{
		Address: "elys17xvmkcpsrtl77yvdl8zh74sw6njg66jdp3tra0",
		Amount:  math.NewInt(95500000),
		Claimed: false,
	},
	{
		Address: "elys1dumrg94tq5j880jnwrgr6x7zfmeudrnahueen9",
		Amount:  math.NewInt(95500000),
		Claimed: false,
	},
	{
		Address: "elys1rzv558svt69c20rx7ez9uzazhs7k0kewrtrcdn",
		Amount:  math.NewInt(95500000),
		Claimed: false,
	},
	{
		Address: "elys1ukdrzvtmg3n0k7p3gan7drj7ym94pu75vmfs7q",
		Amount:  math.NewInt(95400000),
		Claimed: false,
	},
	{
		Address: "elys1um0y6s8nncuzn2pze7q0duhuw90fvgn4rd0zn2",
		Amount:  math.NewInt(95400000),
		Claimed: false,
	},
	{
		Address: "elys1q0rn9njfl5a4c47ruymrxgquwfrwn9jqucljx0",
		Amount:  math.NewInt(95300000),
		Claimed: false,
	},
	{
		Address: "elys10fn9xzprze62g02zm7g3kxgtumu9udhg7qtjx9",
		Amount:  math.NewInt(95200000),
		Claimed: false,
	},
	{
		Address: "elys125nkhdvfwrt9qfleypcggkdwq2fmfgz5mgm2zv",
		Amount:  math.NewInt(95200000),
		Claimed: false,
	},
	{
		Address: "elys1mas532ep6vgeg0stj44raa9jxwuymw00nv4x2j",
		Amount:  math.NewInt(95200000),
		Claimed: false,
	},
	{
		Address: "elys1ef83gzh9cm8l780lqrxw6np475qvvdq4qqqz8t",
		Amount:  math.NewInt(95100000),
		Claimed: false,
	},
	{
		Address: "elys1fwerjzsey0yen5wjgxny6fl29sf66p6ccmhq9l",
		Amount:  math.NewInt(95100000),
		Claimed: false,
	},
	{
		Address: "elys1seha60swa43auhsg75eqwg5xwj305muzruuzyc",
		Amount:  math.NewInt(94800000),
		Claimed: false,
	},
	{
		Address: "elys1dxmd64nxzw4v62jdnhthjve2py4twr8ptxz9dd",
		Amount:  math.NewInt(94700000),
		Claimed: false,
	},
	{
		Address: "elys1gq3wswr4vl4ltd4djdt9d5numfyq0rke6v6zsk",
		Amount:  math.NewInt(94700000),
		Claimed: false,
	},
	{
		Address: "elys159dggfqnd6nrm6a389k4u0g6z59gf9zyqywjxw",
		Amount:  math.NewInt(94600000),
		Claimed: false,
	},
	{
		Address: "elys15t8uyks0s2rs85qv39qh8hv3tuxjwzztz2u5za",
		Amount:  math.NewInt(94500000),
		Claimed: false,
	},
	{
		Address: "elys17gcuws90tm6rsmg93vlan4r3td88yepkn39h4r",
		Amount:  math.NewInt(94500000),
		Claimed: false,
	},
	{
		Address: "elys1cd8w26caw0kmwqz6a2vlnqsyydfcwrkrzu94rc",
		Amount:  math.NewInt(94500000),
		Claimed: false,
	},
	{
		Address: "elys1dqfx7jw4xdsu5v6rxs9ymkxj8v24qvfqpx0hyh",
		Amount:  math.NewInt(94500000),
		Claimed: false,
	},
	{
		Address: "elys1gwt5vn26sp00csy9j9k4gx3kgpptvnweagpypv",
		Amount:  math.NewInt(94500000),
		Claimed: false,
	},
	{
		Address: "elys1tvdl4e9lka08ngpwhywt3x2zlm3m68t3a6wqtj",
		Amount:  math.NewInt(94500000),
		Claimed: false,
	},
	{
		Address: "elys1xwrp0zdqkxxwp0gqrzxhnq853cpcal9us4js5u",
		Amount:  math.NewInt(94500000),
		Claimed: false,
	},
	{
		Address: "elys1y9gvgqsgl9u00f8xe6x74zdxpqsxgwfus7zay4",
		Amount:  math.NewInt(94500000),
		Claimed: false,
	},
	{
		Address: "elys1zyvgexqseu6qtkaad2czzjvch9ann8scf7ja8j",
		Amount:  math.NewInt(94500000),
		Claimed: false,
	},
	{
		Address: "elys1v6hutyc8m4gd8cctxmfe56jxy7a89wsej4z7lv",
		Amount:  math.NewInt(94300000),
		Claimed: false,
	},
	{
		Address: "elys1rhlucdmyc3dd4ma0t0502g8r32jsmcmghkzkdc",
		Amount:  math.NewInt(94200000),
		Claimed: false,
	},
	{
		Address: "elys1vtgg9vjnn9h8uqyt7zlkzpsmutrths44aj757f",
		Amount:  math.NewInt(94100000),
		Claimed: false,
	},
	{
		Address: "elys1qa9535596nm99tmw2uyzlrzged8084dsnf8xle",
		Amount:  math.NewInt(94000000),
		Claimed: false,
	},
	{
		Address: "elys1q8hsvnpjkpakmzqkdx6s48cug9s7gdkz5etumw",
		Amount:  math.NewInt(93900000),
		Claimed: false,
	},
	{
		Address: "elys1qr2lpf7dj586h2s0fcm7qm8w9fsa64rf5mc4tf",
		Amount:  math.NewInt(93900000),
		Claimed: false,
	},
	{
		Address: "elys1ggkdxkxlr38n3cvtnpjs62e4awxks0d53nxkqe",
		Amount:  math.NewInt(93800000),
		Claimed: false,
	},
	{
		Address: "elys13rvwejykl7vfgtepwm9p2q609020fgu67hvpe5",
		Amount:  math.NewInt(93500000),
		Claimed: false,
	},
	{
		Address: "elys15dw4lclv9kcn2yz2hnywzpzlxlcgxv3rpc6gf3",
		Amount:  math.NewInt(93300000),
		Claimed: false,
	},
	{
		Address: "elys1utsp56kls0n0f3hl99ed3dhpdmxxs65xgql6nn",
		Amount:  math.NewInt(92900000),
		Claimed: false,
	},
	{
		Address: "elys1qjchggr92ehuqg3604knhaukkmjed8t0n6us0d",
		Amount:  math.NewInt(92800000),
		Claimed: false,
	},
	{
		Address: "elys12gtnscu9qxvedv79crqp883g8m0p9vah5dc80h",
		Amount:  math.NewInt(92700000),
		Claimed: false,
	},
	{
		Address: "elys17wsng46pxlx7xpqfa83tvvgeqt66kyd48d5wv7",
		Amount:  math.NewInt(92700000),
		Claimed: false,
	},
	{
		Address: "elys1ew3n4lyxgu3klr2d4k6a00ggn78fv7jxnkhfy2",
		Amount:  math.NewInt(92700000),
		Claimed: false,
	},
	{
		Address: "elys1fzhlytvsnn34gehcunf4vu8axtehkvv2rzdg54",
		Amount:  math.NewInt(92700000),
		Claimed: false,
	},
	{
		Address: "elys1j78epr8ewzpsljj4wrjeplrzwaduc0s2qzz6rk",
		Amount:  math.NewInt(92700000),
		Claimed: false,
	},
	{
		Address: "elys1lq5qx9sjqsnkusldlcvlxmyvre0t3e2fwcdr69",
		Amount:  math.NewInt(92700000),
		Claimed: false,
	},
	{
		Address: "elys1srk84xj3kpylf929xd4736n4wavm643jgvr9vc",
		Amount:  math.NewInt(92700000),
		Claimed: false,
	},
	{
		Address: "elys1v5dmhdxl0vneqv8svze0se4fxq0vsxx3jvja0r",
		Amount:  math.NewInt(92700000),
		Claimed: false,
	},
	{
		Address: "elys1x0xnvd0h3dph8lrp050y4cruutz35apyfe2jp4",
		Amount:  math.NewInt(92700000),
		Claimed: false,
	},
	{
		Address: "elys1yhsddussze2gh8r2ejj66e5dwjucdu5zp4gm9q",
		Amount:  math.NewInt(92700000),
		Claimed: false,
	},
	{
		Address: "elys1zwad78d937fm6n83u7yarmlgfypnhlpecqyv7w",
		Amount:  math.NewInt(92700000),
		Claimed: false,
	},
	{
		Address: "elys1w8sntjc4xn9rwm9w8dyf603g05mnkfakdw6j55",
		Amount:  math.NewInt(92600000),
		Claimed: false,
	},
	{
		Address: "elys1yftyaus3cvquttyz3srrytcnzlhs6p4c5urd7x",
		Amount:  math.NewInt(92500000),
		Claimed: false,
	},
	{
		Address: "elys1fxpkqjjtqe9glr7h7arw3ytcqn78kr5ca4nze4",
		Amount:  math.NewInt(92400000),
		Claimed: false,
	},
	{
		Address: "elys13akc9xjenm4y5f7av8gthzsyec3qhzakva7ytl",
		Amount:  math.NewInt(92300000),
		Claimed: false,
	},
	{
		Address: "elys1458eckhxdzfe6nlcus3n2wtn5t0jtrdylne0lp",
		Amount:  math.NewInt(92300000),
		Claimed: false,
	},
	{
		Address: "elys1wzvgwttnhs7sly4ey6p0z8cv2zzcs67svts6up",
		Amount:  math.NewInt(92300000),
		Claimed: false,
	},
	{
		Address: "elys1fq463ptwchsse2uf4tkn94sz6pe7v6mjskataj",
		Amount:  math.NewInt(91700000),
		Claimed: false,
	},
	{
		Address: "elys1z008pllzexhl9xw6w87fvhfe6036putf72a8qq",
		Amount:  math.NewInt(91700000),
		Claimed: false,
	},
	{
		Address: "elys13sk4dwrcywhqas32qrrav7e9lgc86nff6cayjq",
		Amount:  math.NewInt(91300000),
		Claimed: false,
	},
	{
		Address: "elys1mthlecqn2lfgnjp8ypz0fcu5evw9gpsfka4rm4",
		Amount:  math.NewInt(91200000),
		Claimed: false,
	},
	{
		Address: "elys1ymwgg85shgthy4gsp7qf59jpf0pz4c3rnm4hh5",
		Amount:  math.NewInt(91200000),
		Claimed: false,
	},
	{
		Address: "elys17w8vp0g8vh9kr6efp8092mmvaefr8wx42sczc6",
		Amount:  math.NewInt(91000000),
		Claimed: false,
	},
	{
		Address: "elys1hvcs5nltqsec8rgwzu4wn75wjt40pp2077rn0d",
		Amount:  math.NewInt(91000000),
		Claimed: false,
	},
	{
		Address: "elys123qws0a0kv5eh027udjjr54flk7q9vf6cx237z",
		Amount:  math.NewInt(90700000),
		Claimed: false,
	},
	{
		Address: "elys1vyfmzw9dfw4xhef2368pd6mauuwh2765ccl3dm",
		Amount:  math.NewInt(90700000),
		Claimed: false,
	},
	{
		Address: "elys12judf8x3d2t4l2gwx4w47g4r9n0f2nw8r97k0c",
		Amount:  math.NewInt(90500000),
		Claimed: false,
	},
	{
		Address: "elys1psfq5aj39wtfew50s6mhkjrf2fv8cpgnvr6hsg",
		Amount:  math.NewInt(90500000),
		Claimed: false,
	},
	{
		Address: "elys1u0ek20lefga3dyhp2eedt3qwdu6c3kwjamhp6c",
		Amount:  math.NewInt(89900000),
		Claimed: false,
	},
	{
		Address: "elys1uqk9plq9hj94rwpve4yqkfafn33y0drmcyx9ju",
		Amount:  math.NewInt(89700000),
		Claimed: false,
	},
	{
		Address: "elys1ydktl90zv3hzmgsgwzrh8x3wakrd8qfmlqvvlm",
		Amount:  math.NewInt(89700000),
		Claimed: false,
	},
	{
		Address: "elys1zlvmxp5ckc380jxlfxxgrewx54ku9agecz69l8",
		Amount:  math.NewInt(89700000),
		Claimed: false,
	},
	{
		Address: "elys13wwgl5mtfpqvwjpphkn30kxwyhnnljz7vcrpnq",
		Amount:  math.NewInt(89300000),
		Claimed: false,
	},
	{
		Address: "elys1fa9qgeswzrmrxn33p3k6v4t8tqdvqp0gvtucgw",
		Amount:  math.NewInt(89300000),
		Claimed: false,
	},
	{
		Address: "elys1zwc7al09m9js4ua0ct8dph4y4vuc5hjcwap7zd",
		Amount:  math.NewInt(89300000),
		Claimed: false,
	},
	{
		Address: "elys1s7wv7jglr0qy7tnfhu6pqsjlah84d4svkvxj3v",
		Amount:  math.NewInt(89000000),
		Claimed: false,
	},
	{
		Address: "elys1l0zae9n2dejpk6ays9v4zc8vqmtpnk8f3dvyqa",
		Amount:  math.NewInt(88900000),
		Claimed: false,
	},
	{
		Address: "elys1x5sym5lxahyknvwz02x4qjl5drq286t7s6yk6d",
		Amount:  math.NewInt(88900000),
		Claimed: false,
	},
	{
		Address: "elys1jdhk5sdhj0mc4zsuhtdq2p3z6x3cv9darffk0x",
		Amount:  math.NewInt(88700000),
		Claimed: false,
	},
	{
		Address: "elys132ertfqzykwwld9d9yqj5v3v0vfdaxfh7l6asl",
		Amount:  math.NewInt(88600000),
		Claimed: false,
	},
	{
		Address: "elys1w6anyh8znadf7gj032vmv2lfen53nx6dnlhfrw",
		Amount:  math.NewInt(88600000),
		Claimed: false,
	},
	{
		Address: "elys1w7jezmlqaysdfzgy4nm7vafk8s0fmf5uskdpww",
		Amount:  math.NewInt(88600000),
		Claimed: false,
	},
	{
		Address: "elys1xnwjk3tq84q3vvkgg9jwfxwpn0kdhetuvhdfq3",
		Amount:  math.NewInt(88600000),
		Claimed: false,
	},
	{
		Address: "elys1ylee5xqj9ac92k4cgt3yzxzagc533ar9c07n60",
		Amount:  math.NewInt(88200000),
		Claimed: false,
	},
	{
		Address: "elys1r020k5d5p9sc2ytpa38dt56vup8m42yax4x2hw",
		Amount:  math.NewInt(88100000),
		Claimed: false,
	},
	{
		Address: "elys13mdynf4fu4ghmxnfwv9ghq2selhxvq7hw0qp8u",
		Amount:  math.NewInt(87900000),
		Claimed: false,
	},
	{
		Address: "elys1hqckq2gqyrn4zu543jywnwrdjfhfjvlqj5j3j0",
		Amount:  math.NewInt(87900000),
		Claimed: false,
	},
	{
		Address: "elys1s6e8g3hr9v3p80qu5js6htepctwxghfmy60w9h",
		Amount:  math.NewInt(87800000),
		Claimed: false,
	},
	{
		Address: "elys1vvvy0eg9vy783jwyjx7qnx4x7he4vp20mf994x",
		Amount:  math.NewInt(87800000),
		Claimed: false,
	},
	{
		Address: "elys134jlrkevfgru6htdt0e277x2u4ext76rwaenss",
		Amount:  math.NewInt(87600000),
		Claimed: false,
	},
	{
		Address: "elys1395trnj3pu7qzyuu9y0hwzrns2kcpvzk36ed7z",
		Amount:  math.NewInt(87600000),
		Claimed: false,
	},
	{
		Address: "elys1cp8pynm2k8ljm753fpqj2t5yad9ahnpym6krnx",
		Amount:  math.NewInt(87600000),
		Claimed: false,
	},
	{
		Address: "elys1em9hcu5seqwqr67j27k0tw9t4rqw2whe8w4vw0",
		Amount:  math.NewInt(87600000),
		Claimed: false,
	},
	{
		Address: "elys1kczqg5m6hyma2djulx7mnwe5eutneyuwjvex9k",
		Amount:  math.NewInt(87600000),
		Claimed: false,
	},
	{
		Address: "elys18knrpz32ch2rhgaxznt9mucuh2mjfpalu8nt3l",
		Amount:  math.NewInt(87500000),
		Claimed: false,
	},
	{
		Address: "elys1c9e8c53ck7czzx9c4vy9052a82uulfhksvgah4",
		Amount:  math.NewInt(87400000),
		Claimed: false,
	},
	{
		Address: "elys1v59f5p9644psgdcpltcatjaz7tyxvupv79sgu6",
		Amount:  math.NewInt(87400000),
		Claimed: false,
	},
	{
		Address: "elys1r5m978tlcuuuvmaw6xqdj2k2qpv5myxwz7q4mw",
		Amount:  math.NewInt(87300000),
		Claimed: false,
	},
	{
		Address: "elys1udk2c4lvmuktzty0e0le0p4397epwg34snrvsh",
		Amount:  math.NewInt(87300000),
		Claimed: false,
	},
	{
		Address: "elys1utnr4wz2jgg2tzn8gpf3sdx7gc5yvxme5a9ea9",
		Amount:  math.NewInt(87300000),
		Claimed: false,
	},
	{
		Address: "elys15u3828zd4we2sp8j43skqqm0urfapp9cgdccfq",
		Amount:  math.NewInt(87200000),
		Claimed: false,
	},
	{
		Address: "elys1ad3v43cdfcq762hpvlnhjcgc4z434kpfgjz4q7",
		Amount:  math.NewInt(87200000),
		Claimed: false,
	},
	{
		Address: "elys1q9f64nzazmnezs6dqkv5wdwhn6zckvz3d50gh7",
		Amount:  math.NewInt(87100000),
		Claimed: false,
	},
	{
		Address: "elys1edc37szy2m9a5xmhlaezy0ppn0c9hwt3ccvutx",
		Amount:  math.NewInt(87000000),
		Claimed: false,
	},
	{
		Address: "elys17cjfay6kggl45xqagdv7p3qnqxha93v9synxk7",
		Amount:  math.NewInt(86900000),
		Claimed: false,
	},
	{
		Address: "elys1lm0kg3aex3fqxlk7feuft7dtvaqgfwmvnfmjd8",
		Amount:  math.NewInt(86900000),
		Claimed: false,
	},
	{
		Address: "elys1uak3aga7arn282sf349ck92pf8uadycsqs7nn3",
		Amount:  math.NewInt(86700000),
		Claimed: false,
	},
	{
		Address: "elys1tjptck8amxsq7uwqfytrna7zamzh2yz98fdy4l",
		Amount:  math.NewInt(86600000),
		Claimed: false,
	},
	{
		Address: "elys1xxue5c80857ugxg0wjqx2mweggpphcc8y280pj",
		Amount:  math.NewInt(86600000),
		Claimed: false,
	},
	{
		Address: "elys14jhzh029gd9fc5z3y4elczlsjxz9a9fnu5g8wr",
		Amount:  math.NewInt(86500000),
		Claimed: false,
	},
	{
		Address: "elys1tcrjwejyhgm6z7ll856zy0f4sy3y367yqgcmm0",
		Amount:  math.NewInt(86500000),
		Claimed: false,
	},
	{
		Address: "elys1mx92dj2nwd49cs6dfqgjw48muzcplgetcgzgth",
		Amount:  math.NewInt(86400000),
		Claimed: false,
	},
	{
		Address: "elys12g4wk37tm65fjgvdtn2mhmy3c5ak5rf62g08d7",
		Amount:  math.NewInt(86200000),
		Claimed: false,
	},
	{
		Address: "elys1nqhqkpsdpng00nrlr4a3rx5emyqetkmkyjjrmm",
		Amount:  math.NewInt(86200000),
		Claimed: false,
	},
	{
		Address: "elys10yxjeppxhr4yae2jp6tl6ephlt0qjnm5fzmc65",
		Amount:  math.NewInt(85900000),
		Claimed: false,
	},
	{
		Address: "elys133zkt3y0mcetla2mdknnmalgrm7c9ln5addhkd",
		Amount:  math.NewInt(85900000),
		Claimed: false,
	},
	{
		Address: "elys176y260mk65ekz4cxlju9njtfauxnat5e6zv20l",
		Amount:  math.NewInt(85900000),
		Claimed: false,
	},
	{
		Address: "elys177vfy23f3m9e6k5umfwy8dsp6efqcd68gty9wc",
		Amount:  math.NewInt(85900000),
		Claimed: false,
	},
	{
		Address: "elys18482ap5dzwxknrlxl4f50d0l6em3afcs3cg5l7",
		Amount:  math.NewInt(85900000),
		Claimed: false,
	},
	{
		Address: "elys1d74f6c7rgyyym57530rc2tt88fnl6wpdvavqdw",
		Amount:  math.NewInt(85900000),
		Claimed: false,
	},
	{
		Address: "elys1damt355s4kvqtcf4w56p89v5krxlyx5w9tdjal",
		Amount:  math.NewInt(85900000),
		Claimed: false,
	},
	{
		Address: "elys1g8y2eeqp54fe9n6lfwvrqjln3hxg9l355sqrzu",
		Amount:  math.NewInt(85900000),
		Claimed: false,
	},
	{
		Address: "elys1pe6utrx5rm90nerxx2qjl40q9nfh9z6a353w6y",
		Amount:  math.NewInt(85900000),
		Claimed: false,
	},
	{
		Address: "elys1x2s5ccf9mgcr43ywxsc95reae9va58ygz5ddq3",
		Amount:  math.NewInt(85900000),
		Claimed: false,
	},
	{
		Address: "elys1x79yqj836j0jxsx9s0u4c5t0vwq88jayv2uzg6",
		Amount:  math.NewInt(85900000),
		Claimed: false,
	},
	{
		Address: "elys12ghctmu7cmywr3uyxqpq44t247uxlhr84mxy4y",
		Amount:  math.NewInt(85800000),
		Claimed: false,
	},
	{
		Address: "elys140799x36062x05rdln3hpjht8fvhtugaxz98uh",
		Amount:  math.NewInt(85800000),
		Claimed: false,
	},
	{
		Address: "elys15aj2qkj6gqkggjjqgdmncyvw87dpzakcfyvd2e",
		Amount:  math.NewInt(85700000),
		Claimed: false,
	},
	{
		Address: "elys1xuhe0fvy56um4zqmvfa75xm6ldf8xqf3kqe7z4",
		Amount:  math.NewInt(85700000),
		Claimed: false,
	},
	{
		Address: "elys1hxxhd2esyr7rch2c72fj8tal86rj7kn0tam89y",
		Amount:  math.NewInt(85600000),
		Claimed: false,
	},
	{
		Address: "elys1jlhzpfrq8mqq8ujnefctekj037k0ke232qp4e6",
		Amount:  math.NewInt(85500000),
		Claimed: false,
	},
	{
		Address: "elys1cx3vk0lprfvvyfagkpdc5jq88curh0e25f3eny",
		Amount:  math.NewInt(85400000),
		Claimed: false,
	},
	{
		Address: "elys14wymarwujxvq0qk5twxl3cnqad72802z8sg3un",
		Amount:  math.NewInt(85300000),
		Claimed: false,
	},
	{
		Address: "elys1jdven2lq8nk6n5qzy2xnrylnf4kzwyrm0s2nuk",
		Amount:  math.NewInt(85300000),
		Claimed: false,
	},
	{
		Address: "elys1n8lrqmtjveq6kjw5tkvqyf3hpvg3g7hnv03hl0",
		Amount:  math.NewInt(85300000),
		Claimed: false,
	},
	{
		Address: "elys1w7ctwaysyafq7kvzfad0vj7tr3gczfl7d83pfd",
		Amount:  math.NewInt(85200000),
		Claimed: false,
	},
	{
		Address: "elys106vxud3duq5flrskezxrkyanu9jx98hkm99zlj",
		Amount:  math.NewInt(85100000),
		Claimed: false,
	},
	{
		Address: "elys12qf5j8dl6pqtrepkun4jalpejdmxrw38u36c5k",
		Amount:  math.NewInt(85100000),
		Claimed: false,
	},
	{
		Address: "elys1czj2v58w9kl6k2rw64lkdkcntt08urc7a76kdz",
		Amount:  math.NewInt(85100000),
		Claimed: false,
	},
	{
		Address: "elys10fxgaz8zv3wshzcgp222s823fkj8m8kpaq8vxq",
		Amount:  math.NewInt(85000000),
		Claimed: false,
	},
	{
		Address: "elys1q70cr3avhsqwyrklmdmfeg9mxngv88587rrdqt",
		Amount:  math.NewInt(84900000),
		Claimed: false,
	},
	{
		Address: "elys18ku89y82j346q3g5la80e0qgf27gl3udc2txlv",
		Amount:  math.NewInt(84700000),
		Claimed: false,
	},
	{
		Address: "elys1aqy852guxqr0e93zkgxv9ftx55ypa4d6rwmlur",
		Amount:  math.NewInt(84700000),
		Claimed: false,
	},
	{
		Address: "elys1zj00yfhupxrz56p727x8lwzhccapd6w4sgw3mr",
		Amount:  math.NewInt(84700000),
		Claimed: false,
	},
	{
		Address: "elys19cnmuuwue83yxs4z46ph4ldt309fnru98cpj9w",
		Amount:  math.NewInt(84600000),
		Claimed: false,
	},
	{
		Address: "elys188sxst5tpvlp2hrheul9r2q87f0ggevu5zsw2h",
		Amount:  math.NewInt(84500000),
		Claimed: false,
	},
	{
		Address: "elys10m8hmhvx42js2dupem6qspzhegcl5aw4vwgtyp",
		Amount:  math.NewInt(84400000),
		Claimed: false,
	},
	{
		Address: "elys1czd9rr80pp2rtht9akdhk542e4qyktqs4yzv7g",
		Amount:  math.NewInt(84400000),
		Claimed: false,
	},
	{
		Address: "elys142cl70a6rsp56x7myfz67g7m0f5x84w574vpnr",
		Amount:  math.NewInt(84100000),
		Claimed: false,
	},
	{
		Address: "elys1qgdkdcraut5gkvtzv4yfmmytedze9979a89k3r",
		Amount:  math.NewInt(84100000),
		Claimed: false,
	},
	{
		Address: "elys1utf3kwuraecdq8zlhm49gg2pknpuvxhcez75w6",
		Amount:  math.NewInt(84100000),
		Claimed: false,
	},
	{
		Address: "elys10jn6u0tasjkjgzrxkvxj8v83kkekaqkegrd77p",
		Amount:  math.NewInt(84000000),
		Claimed: false,
	},
	{
		Address: "elys17uu57q23s8mtmtt8tru09y9qytudxpsj8dufs0",
		Amount:  math.NewInt(83900000),
		Claimed: false,
	},
	{
		Address: "elys1a4hran6y2znq8ppe0gan25yyc2smq286van26l",
		Amount:  math.NewInt(83900000),
		Claimed: false,
	},
	{
		Address: "elys1jh0k89k2ea54vrdcwrzc7sketz27d8x526aefz",
		Amount:  math.NewInt(83900000),
		Claimed: false,
	},
	{
		Address: "elys1aacda22qyffsmthshzyc8kmvs39jz7rfk72x5a",
		Amount:  math.NewInt(83700000),
		Claimed: false,
	},
	{
		Address: "elys1t52wekkse83q33ea6d5x8gftkrzymperx4shsr",
		Amount:  math.NewInt(83600000),
		Claimed: false,
	},
	{
		Address: "elys1ykkc4waewch0n4zhz5rap2tkhv38qjec32hwau",
		Amount:  math.NewInt(83500000),
		Claimed: false,
	},
	{
		Address: "elys1ndhk22uy66ttzydggc32mm67l536zkgd0cg8tm",
		Amount:  math.NewInt(83300000),
		Claimed: false,
	},
	{
		Address: "elys16wt0ujqeuy0z78dd63a63gshftqnv2z2r30s9a",
		Amount:  math.NewInt(83100000),
		Claimed: false,
	},
	{
		Address: "elys1lfmh79xjmsjgatnucq9pl72uzc7hzfu57vm72q",
		Amount:  math.NewInt(82900000),
		Claimed: false,
	},
	{
		Address: "elys1ljrr3af88xhs7c39946syw35gx96rexk4wuutv",
		Amount:  math.NewInt(82900000),
		Claimed: false,
	},
	{
		Address: "elys1yqccxuht57ne4me0u2h396a6mwdqza0ydnmhcc",
		Amount:  math.NewInt(82900000),
		Claimed: false,
	},
	{
		Address: "elys1p5cpd6a6dh3gwgv4cl72jhceeac40pp2gvpjwx",
		Amount:  math.NewInt(82800000),
		Claimed: false,
	},
	{
		Address: "elys19wqrehrqduk0djz7ud6xw2u07e5kevh6gj3p6v",
		Amount:  math.NewInt(82700000),
		Claimed: false,
	},
	{
		Address: "elys1xck8mcvd2alf8cyt7qwym5myh0lzkv3muwj974",
		Amount:  math.NewInt(82500000),
		Claimed: false,
	},
	{
		Address: "elys15enkvzflchjrjnshpfvgus9uzf8qxc7ar6j8hr",
		Amount:  math.NewInt(82400000),
		Claimed: false,
	},
	{
		Address: "elys16t0aun0w36uyfhz3asjlrcy7a7n6n9vg9p3pyu",
		Amount:  math.NewInt(82400000),
		Claimed: false,
	},
	{
		Address: "elys1kk9zap26ymqppjyy99w04vxsurvx08z6h3te7s",
		Amount:  math.NewInt(82300000),
		Claimed: false,
	},
	{
		Address: "elys1tyhe2aa89kr7gmrckumc8msjrhg7c9x72n0dkf",
		Amount:  math.NewInt(82000000),
		Claimed: false,
	},
	{
		Address: "elys1nre4cf6am4qt2gsxru0m89mkz7d04gd4sfkn59",
		Amount:  math.NewInt(81900000),
		Claimed: false,
	},
	{
		Address: "elys159f2g2zn00n8qmru8xzzlqeax3mft0mgme8n4q",
		Amount:  math.NewInt(81700000),
		Claimed: false,
	},
	{
		Address: "elys109xtyezdwyd5a5n7px4ulm3ugdxt9xk6z478p3",
		Amount:  math.NewInt(81600000),
		Claimed: false,
	},
	{
		Address: "elys15x3c8z6uu5vcwj5p9z8ljxrvy7e3gwg0aryskc",
		Amount:  math.NewInt(81600000),
		Claimed: false,
	},
	{
		Address: "elys166y5n23qwl4u360fduq3zjqps9graetnmurw27",
		Amount:  math.NewInt(81600000),
		Claimed: false,
	},
	{
		Address: "elys1l87j8xyvh3cn6dfvgrx3t09jupe5c707jwny4d",
		Amount:  math.NewInt(81600000),
		Claimed: false,
	},
	{
		Address: "elys1p4cv3gslg593x0jetxhe0fkry2c5gj36zcq0px",
		Amount:  math.NewInt(81600000),
		Claimed: false,
	},
	{
		Address: "elys1rm9k7l7hvdq4tg6gp2xdaz9h2a0kktmddfxdhj",
		Amount:  math.NewInt(81600000),
		Claimed: false,
	},
	{
		Address: "elys1w5vp4q4zw7fsnhlj3y6dd8vknm7ddm3ek3rsql",
		Amount:  math.NewInt(81600000),
		Claimed: false,
	},
	{
		Address: "elys1wg9cpmx9tfjxg5chrtlg9nlgq3gdlv968swgqe",
		Amount:  math.NewInt(81600000),
		Claimed: false,
	},
	{
		Address: "elys1wup3xsnwszfyaec6tavjy4sy0l55x5a528rxpw",
		Amount:  math.NewInt(81600000),
		Claimed: false,
	},
	{
		Address: "elys1xzrez403aukll7a2v3sdcv06phsj35qhr8xnzm",
		Amount:  math.NewInt(81600000),
		Claimed: false,
	},
	{
		Address: "elys1yxdyuvncwhdrm2nw9z35cpel05jd3aamxpa7nx",
		Amount:  math.NewInt(81600000),
		Claimed: false,
	},
	{
		Address: "elys10app9zaajt6m2keng2rfp5l2kydv7wymfmjewl",
		Amount:  math.NewInt(81500000),
		Claimed: false,
	},
	{
		Address: "elys1uesgtap0h5gllq6lx9626mfxmkkg60466k7q7x",
		Amount:  math.NewInt(81500000),
		Claimed: false,
	},
	{
		Address: "elys1vkw0lm7d5hyctguslkx3tjnkssxl752pa6t3jy",
		Amount:  math.NewInt(81500000),
		Claimed: false,
	},
	{
		Address: "elys1638ya9nn05tjghk7ewyv6epwfxzzx2wtcp0t4y",
		Amount:  math.NewInt(81400000),
		Claimed: false,
	},
	{
		Address: "elys17f5trsf9z6tyfnadfzrxgsk8wcs9hyealp8ce8",
		Amount:  math.NewInt(81400000),
		Claimed: false,
	},
	{
		Address: "elys19pvf97ctj0j4utdcagkzams03lgm7gmwf5rw0j",
		Amount:  math.NewInt(81400000),
		Claimed: false,
	},
	{
		Address: "elys1exzxvurclqxhgvpvvlz3lqga7u49qpxczfxeqx",
		Amount:  math.NewInt(81400000),
		Claimed: false,
	},
	{
		Address: "elys1vlr5unhjumn2muk6ajkz0e7gdelp2jc6dgfnm7",
		Amount:  math.NewInt(81400000),
		Claimed: false,
	},
	{
		Address: "elys1wpvg7lae6j835q5hfqd0408s2eca5mhwfgsr86",
		Amount:  math.NewInt(81200000),
		Claimed: false,
	},
	{
		Address: "elys16xd0rzp26eldftgtz74sxkq2qssjewh0cak3z2",
		Amount:  math.NewInt(81100000),
		Claimed: false,
	},
	{
		Address: "elys1896ehs324vec5cl033r8j4r3ygfrh444gvvyrc",
		Amount:  math.NewInt(81100000),
		Claimed: false,
	},
	{
		Address: "elys1crp60ua6z32645kkqy9m3exgemrhmfheahhjzu",
		Amount:  math.NewInt(81100000),
		Claimed: false,
	},
	{
		Address: "elys1wgr69s3f9p6m47yhntyup2jaxa8a3ktg4lpyzk",
		Amount:  math.NewInt(80900000),
		Claimed: false,
	},
	{
		Address: "elys1x53u7kqxrq6mzvl4r26vvfnf5un5kpnrs3hxeq",
		Amount:  math.NewInt(80900000),
		Claimed: false,
	},
	{
		Address: "elys12udm09e75jl7jntamptsgs0yt60prvtpre0acy",
		Amount:  math.NewInt(80800000),
		Claimed: false,
	},
	{
		Address: "elys12ewtvn46crn2sre8pqv3zgk949nxhyxepkkzpq",
		Amount:  math.NewInt(80700000),
		Claimed: false,
	},
	{
		Address: "elys19p3ep5l8nyue75q0j0jpcyshqet6ph5nuy4qph",
		Amount:  math.NewInt(80700000),
		Claimed: false,
	},
	{
		Address: "elys1z6x9qmdnanq5dndcu6576dpdq48xsavz5dl8kp",
		Amount:  math.NewInt(80700000),
		Claimed: false,
	},
	{
		Address: "elys1x8w82rmam0nnmumqqwr02lezsafnm28rrt885t",
		Amount:  math.NewInt(80500000),
		Claimed: false,
	},
	{
		Address: "elys12mjxus7s5crvehpkjyzq0xtxg2adnl3vqvzf4s",
		Amount:  math.NewInt(80400000),
		Claimed: false,
	},
	{
		Address: "elys1za3wczlyeulswx73ccu6edny7l2hzyw5p04a3g",
		Amount:  math.NewInt(80300000),
		Claimed: false,
	},
	{
		Address: "elys18qdggalvs9qa6l9vmvxxgef8hwukshrpzxmlc9",
		Amount:  math.NewInt(80200000),
		Claimed: false,
	},
	{
		Address: "elys1ndpus8r45qetljv6dxcv526srfmmknc9d9evkc",
		Amount:  math.NewInt(80200000),
		Claimed: false,
	},
	{
		Address: "elys1dkjr3vey927m9r8cnlteg5cfzxqjxgxl6asrzj",
		Amount:  math.NewInt(80100000),
		Claimed: false,
	},
	{
		Address: "elys195xp9hd0zj3mmfag3shnj24686hf9twvwmw94c",
		Amount:  math.NewInt(80000000),
		Claimed: false,
	},
	{
		Address: "elys1nu5h45l2zpz56rludfwyhlpnztgd55e5z3uyww",
		Amount:  math.NewInt(80000000),
		Claimed: false,
	},
	{
		Address: "elys1z9lask2kegtt25l2vs9qm80p4tfwzwptc6rlax",
		Amount:  math.NewInt(79900000),
		Claimed: false,
	},
	{
		Address: "elys127y2syu9cpfc78kstgwmvag5t246xm8r2x73fq",
		Amount:  math.NewInt(79800000),
		Claimed: false,
	},
	{
		Address: "elys1ttwqxg65lqd2ak968enjmfz6fm4c7gyp2x2uzy",
		Amount:  math.NewInt(79800000),
		Claimed: false,
	},
	{
		Address: "elys1yrgrvmawrtvyzf7hhu2m3qchqm80k2glknzj8k",
		Amount:  math.NewInt(79700000),
		Claimed: false,
	},
	{
		Address: "elys15xkj85umyxj0v6e99pk7puuk7cdvhgzx3f3s3t",
		Amount:  math.NewInt(79600000),
		Claimed: false,
	},
	{
		Address: "elys1y0rjc69084zw0jchqmv7rz664j8kjdxks9p9ce",
		Amount:  math.NewInt(79600000),
		Claimed: false,
	},
	{
		Address: "elys15a3a7lszqezajfal876h5uaad8awtrcnuskqew",
		Amount:  math.NewInt(79500000),
		Claimed: false,
	},
	{
		Address: "elys1a95flnvk0wxv6drvznrz2swx8ex3clz0jzwffy",
		Amount:  math.NewInt(79300000),
		Claimed: false,
	},
	{
		Address: "elys1vfk3rrh048fvz82jjxkkcxx709u87q8kl9a8tv",
		Amount:  math.NewInt(79300000),
		Claimed: false,
	},
	{
		Address: "elys1czsp2rq0n7m0z35404g9pntqk2h0fduk6cq9hj",
		Amount:  math.NewInt(79200000),
		Claimed: false,
	},
	{
		Address: "elys1heqte6tv99kfestpjnkjevrksdecdeu9vajvc7",
		Amount:  math.NewInt(79200000),
		Claimed: false,
	},
	{
		Address: "elys175nf2a4cc9mjdmq2jayyjmvqs7w3hjp7dqsduc",
		Amount:  math.NewInt(79100000),
		Claimed: false,
	},
	{
		Address: "elys1vavuhv6x84680z7drkvwsf79m9wnl2vh6cjegq",
		Amount:  math.NewInt(79100000),
		Claimed: false,
	},
	{
		Address: "elys1e590j6dt246gxjv6qx94yc4qvggm93t0xs0uns",
		Amount:  math.NewInt(78900000),
		Claimed: false,
	},
	{
		Address: "elys1rqg2m7gek8avfptcj23yg825khl7vdvhmdjlqs",
		Amount:  math.NewInt(78800000),
		Claimed: false,
	},
	{
		Address: "elys1ctjrkqfways3ry70zf5kq9aug5aegxsevuzqxp",
		Amount:  math.NewInt(78700000),
		Claimed: false,
	},
	{
		Address: "elys1x08mzqa84vdfdmyzwfydl0d0wagzp86yjgtkql",
		Amount:  math.NewInt(78700000),
		Claimed: false,
	},
	{
		Address: "elys1a6qkrtflmetqf3gh0mvettq4lhhe8f94sxlmmh",
		Amount:  math.NewInt(78200000),
		Claimed: false,
	},
	{
		Address: "elys1nqzlfr4lk00lhyvm4x2x880q2l0de8tsn85p3p",
		Amount:  math.NewInt(78100000),
		Claimed: false,
	},
	{
		Address: "elys1nxzn44xz349gyazqr42a3u5pyvtffqg5pna3tt",
		Amount:  math.NewInt(77800000),
		Claimed: false,
	},
	{
		Address: "elys16f277w06q5xvfu8j6lm9eg8cztg8nc7m2vawrz",
		Amount:  math.NewInt(77700000),
		Claimed: false,
	},
	{
		Address: "elys1g44c5z8u7g25s68ats2r23hwfmta6p2hepl390",
		Amount:  math.NewInt(77600000),
		Claimed: false,
	},
	{
		Address: "elys1nugzpq4yyr54jp9gr3nna4ezys9syhax4dsp3g",
		Amount:  math.NewInt(77500000),
		Claimed: false,
	},
	{
		Address: "elys1fsehy5vna2mrud6jry4tpz56etqunaphfwu23t",
		Amount:  math.NewInt(77400000),
		Claimed: false,
	},
	{
		Address: "elys1vcqgp3qf9wa8dkng3ys4mn9azl7te8ldh25x00",
		Amount:  math.NewInt(77400000),
		Claimed: false,
	},
	{
		Address: "elys1y2cmk66esng3qehzhdvwnhns5wlm7pe744danj",
		Amount:  math.NewInt(77400000),
		Claimed: false,
	},
	{
		Address: "elys16h0xdhylq4ph4jdfj2ukxwhmkz06xqznhuxnys",
		Amount:  math.NewInt(77300000),
		Claimed: false,
	},
	{
		Address: "elys1rk6qmhyk43pe7nwa55xj3ge5cndc38m4gzsd2e",
		Amount:  math.NewInt(77300000),
		Claimed: false,
	},
	{
		Address: "elys128pk7d3mjw4nnhffhecwn4ga5lu0emlj5qftpj",
		Amount:  math.NewInt(77100000),
		Claimed: false,
	},
	{
		Address: "elys1r5muk8te0z2url8kntgjst7wygnfmr8xyvk6hh",
		Amount:  math.NewInt(77100000),
		Claimed: false,
	},
	{
		Address: "elys19ymutk8t0x0w5x6z9l4zqwms5esss2rc4yj3zf",
		Amount:  math.NewInt(77000000),
		Claimed: false,
	},
	{
		Address: "elys1czse5lpr2606nfw329w7jlvqlp8uqyj0j2adv4",
		Amount:  math.NewInt(77000000),
		Claimed: false,
	},
	{
		Address: "elys143zr7lfdtqj85v8gsrum4fjplmpnuygdgvx44v",
		Amount:  math.NewInt(76700000),
		Claimed: false,
	},
	{
		Address: "elys1sm2wxa7wheczfpf7qpkc3jy74gl3wjksw7ggaf",
		Amount:  math.NewInt(76600000),
		Claimed: false,
	},
	{
		Address: "elys1z7w9kzvhxtdupczykq3vgtggpdjhtndfzhjt80",
		Amount:  math.NewInt(76600000),
		Claimed: false,
	},
	{
		Address: "elys1444q68trhcxm6f8p3s63txt8l6qt5alswanemv",
		Amount:  math.NewInt(76500000),
		Claimed: false,
	},
	{
		Address: "elys1rxzzng878lkrzhrs8jw3fj4l0flpd3wz2gajdq",
		Amount:  math.NewInt(76300000),
		Claimed: false,
	},
	{
		Address: "elys1fymu6xy9z6eupe4g5unqh36dnltxg9q8gegghc",
		Amount:  math.NewInt(76200000),
		Claimed: false,
	},
	{
		Address: "elys1askstu6vgmsund8z8vtadpaxv2n8vanvthvnkh",
		Amount:  math.NewInt(76100000),
		Claimed: false,
	},
	{
		Address: "elys1kjlrmfnqeufvmjrcqs88jjye6sgkk8vpg7y8vt",
		Amount:  math.NewInt(76000000),
		Claimed: false,
	},
	{
		Address: "elys165nfl0f9cc7up2s6k8cfqqvm47mnr2d5ydeu5c",
		Amount:  math.NewInt(75900000),
		Claimed: false,
	},
	{
		Address: "elys1c08tg5mgqp5u4ynvm38ap3msz7hrgvuzw45sjn",
		Amount:  math.NewInt(75900000),
		Claimed: false,
	},
	{
		Address: "elys1e0v5qu9hknmtpy3ha0ye9vh2pqvsdn8hlqcyn3",
		Amount:  math.NewInt(75800000),
		Claimed: false,
	},
	{
		Address: "elys129urhm6ektq84xt7mluzk7tddev5upexzy5g2w",
		Amount:  math.NewInt(75700000),
		Claimed: false,
	},
	{
		Address: "elys1es4gclnpjjlqppjvn5e4hknehm6sm494uv7qs5",
		Amount:  math.NewInt(75700000),
		Claimed: false,
	},
	{
		Address: "elys1r09ld684hdfc0jcvzj3p09sh36hvz3wy099s7g",
		Amount:  math.NewInt(75700000),
		Claimed: false,
	},
	{
		Address: "elys188ky3gqx4safz602q9qalwjnv0p0xfwm8x904a",
		Amount:  math.NewInt(75600000),
		Claimed: false,
	},
	{
		Address: "elys1u9uhljnu555pnamtak0amue3yv7z7tvf3eeqnx",
		Amount:  math.NewInt(75600000),
		Claimed: false,
	},
	{
		Address: "elys1p7wxkpwetplhz82ye7mjt7txzyzd8vssyts0f0",
		Amount:  math.NewInt(75400000),
		Claimed: false,
	},
	{
		Address: "elys1fsd353y8lewe4hhc5trpltrctsjsesh67nfdvp",
		Amount:  math.NewInt(75300000),
		Claimed: false,
	},
	{
		Address: "elys14z94ewufvznwsvajcx7l9dxjn5sgvkqjxnaute",
		Amount:  math.NewInt(75100000),
		Claimed: false,
	},
	{
		Address: "elys1cfcawzd43lxefjw4gcfau82ndlrcsm2cecfajv",
		Amount:  math.NewInt(75100000),
		Claimed: false,
	},
	{
		Address: "elys1dy099nfruh6rc3um6lpxjfxuthm3n8pultct6u",
		Amount:  math.NewInt(75000000),
		Claimed: false,
	},
	{
		Address: "elys1my56tj9lrtsj0lelfwyhw4k8vree357rf65z7s",
		Amount:  math.NewInt(75000000),
		Claimed: false,
	},
	{
		Address: "elys1v2ue9n0vq2pesw2yv0gc5rn5efwc6wrk4kyq7w",
		Amount:  math.NewInt(75000000),
		Claimed: false,
	},
	{
		Address: "elys1cy7h5znrg6uwqryutvdc0l9zgndvasj44hw840",
		Amount:  math.NewInt(74900000),
		Claimed: false,
	},
	{
		Address: "elys1qcwksa3xp3yqfkeztsycvftkvqwr0pl4u9v7mq",
		Amount:  math.NewInt(74900000),
		Claimed: false,
	},
	{
		Address: "elys1t9qrpa8jcnjtxrn00pr59vm7lqylqs4h4fzhvt",
		Amount:  math.NewInt(74700000),
		Claimed: false,
	},
	{
		Address: "elys1n2amtckvcwxte5feckvjlwt3shpn06jsqyqjz4",
		Amount:  math.NewInt(74500000),
		Claimed: false,
	},
	{
		Address: "elys1nmfscw5nled3vtw4khk80mhw9uxhssx7mkfgyy",
		Amount:  math.NewInt(74400000),
		Claimed: false,
	},
	{
		Address: "elys1tzfxftmj0espy0ydfmh4uv7u0pngyvq9kcd2rk",
		Amount:  math.NewInt(74400000),
		Claimed: false,
	},
	{
		Address: "elys1rnjcrw9kxl5phzf45eglsm6g8skpk4jg8jq54l",
		Amount:  math.NewInt(74300000),
		Claimed: false,
	},
	{
		Address: "elys12mdqnpglrs33j9m08pj0kykrdzx889hwq6s876",
		Amount:  math.NewInt(74200000),
		Claimed: false,
	},
	{
		Address: "elys1xffm6wzyyq092kfvf7ng77jjr90tagtz6vx5uu",
		Amount:  math.NewInt(74200000),
		Claimed: false,
	},
	{
		Address: "elys1rf8rnk935hrth9snk3agzeqmrswc9rp4vqgh32",
		Amount:  math.NewInt(74100000),
		Claimed: false,
	},
	{
		Address: "elys18x7hm29weyew6slna66rjpf5l3f5hvcak8jxux",
		Amount:  math.NewInt(74000000),
		Claimed: false,
	},
	{
		Address: "elys1nyh50ga0uvmad4v6c548gp2dtav74rv43kl8t2",
		Amount:  math.NewInt(74000000),
		Claimed: false,
	},
	{
		Address: "elys1c7hkkvqaxwshjtxhfml5deurapw8mwstycrr49",
		Amount:  math.NewInt(73900000),
		Claimed: false,
	},
	{
		Address: "elys1j67dh7jpm9wq7jretplta4pyas5w54f7pmrnu7",
		Amount:  math.NewInt(73900000),
		Claimed: false,
	},
	{
		Address: "elys1wyyue7q4z2dqq8v9vr7qhlsjrqnglunzczpqre",
		Amount:  math.NewInt(73800000),
		Claimed: false,
	},
	{
		Address: "elys19ja83y6qcppsku4krk86lawqh6q702vrx8ka2e",
		Amount:  math.NewInt(73600000),
		Claimed: false,
	},
	{
		Address: "elys1ltjgcxepqwgm75wnjmpm6dz5rteelhyd57ts8j",
		Amount:  math.NewInt(73500000),
		Claimed: false,
	},
	{
		Address: "elys1t7ftpldze6cthukkgaz85nzm7xahh4fauvtxtf",
		Amount:  math.NewInt(73500000),
		Claimed: false,
	},
	{
		Address: "elys1vszv0d8wvc5jqju0wdg63v5y0jq4lh9k6xy3qv",
		Amount:  math.NewInt(73400000),
		Claimed: false,
	},
	{
		Address: "elys146y7j2k96unezwpsclj40q3ncsyxcgwysq4xfv",
		Amount:  math.NewInt(73300000),
		Claimed: false,
	},
	{
		Address: "elys1j3mtmxee0tupnqgp4ep4e6kdtpdf8s7nnhfxzy",
		Amount:  math.NewInt(73300000),
		Claimed: false,
	},
	{
		Address: "elys1p7hymf663g4wlq63xspdcraj7cawnvkzvuq0td",
		Amount:  math.NewInt(73300000),
		Claimed: false,
	},
	{
		Address: "elys1q7uxyjs9vkesn469nyujuvjl6k70usp72kry7g",
		Amount:  math.NewInt(73200000),
		Claimed: false,
	},
	{
		Address: "elys1xdffefsxuyrcjmjtc59zru9v09vs6zqmvx77nl",
		Amount:  math.NewInt(73200000),
		Claimed: false,
	},
	{
		Address: "elys1av9aejgnpv6js5smdwzle0372s5wwwsmstchyy",
		Amount:  math.NewInt(73100000),
		Claimed: false,
	},
	{
		Address: "elys1dw20gphswk0cy82w6pvvhx7k7d22wysvej2e6d",
		Amount:  math.NewInt(73100000),
		Claimed: false,
	},
	{
		Address: "elys1q9lnppqmm8y5uryrxad2hy7zreq9xhv4p3a6ya",
		Amount:  math.NewInt(73100000),
		Claimed: false,
	},
	{
		Address: "elys1y4zzkaprt9cxgfurajlt5m4wvycec7ahk5u7z8",
		Amount:  math.NewInt(73100000),
		Claimed: false,
	},
	{
		Address: "elys1cqvzcymx33cytuuk408mujaqducj6474eay0wr",
		Amount:  math.NewInt(73000000),
		Claimed: false,
	},
	{
		Address: "elys17jsxgde8jcqs3pdeejfazqj5mynknc00tptutz",
		Amount:  math.NewInt(72800000),
		Claimed: false,
	},
	{
		Address: "elys1aeq884fc02c68n27hgwvdghn4v2jt0sl028lql",
		Amount:  math.NewInt(72800000),
		Claimed: false,
	},
	{
		Address: "elys1tlqwkl5akz54v0042f6mw90zvys97da089wduj",
		Amount:  math.NewInt(72800000),
		Claimed: false,
	},
	{
		Address: "elys1exm2hg6jdndcqe54nyxpmqpmshq9tmtf8k2jh2",
		Amount:  math.NewInt(72700000),
		Claimed: false,
	},
	{
		Address: "elys1sc8fsdtevxvjx8gaazwz79rjshqj7g42mx4m2m",
		Amount:  math.NewInt(72700000),
		Claimed: false,
	},
	{
		Address: "elys1cnsuyhej0f2nuaqjk37tdjng39pkd5whpj9l7p",
		Amount:  math.NewInt(72500000),
		Claimed: false,
	},
	{
		Address: "elys1dfwa5lk5k6dl0sqffnuqs8rasmu8ydujhakptj",
		Amount:  math.NewInt(72500000),
		Claimed: false,
	},
	{
		Address: "elys1l8dcdq6c8p64mkz20ul7wev2mz5xazkuqmwa26",
		Amount:  math.NewInt(72500000),
		Claimed: false,
	},
	{
		Address: "elys1wdunel5fkv9dsp5xm2hvuwejuy97xg445s4jmk",
		Amount:  math.NewInt(72200000),
		Claimed: false,
	},
	{
		Address: "elys10935e29j9kcdjgt5vg4kf80z6lkchwx43kad3r",
		Amount:  math.NewInt(72100000),
		Claimed: false,
	},
	{
		Address: "elys1e4rp9hrhxgw7gdtntvhrus2s0j6k80jlpy27ng",
		Amount:  math.NewInt(72100000),
		Claimed: false,
	},
	{
		Address: "elys1hgjn0uvg8c0utd7y946tu37gvwpdwqnff2x062",
		Amount:  math.NewInt(72100000),
		Claimed: false,
	},
	{
		Address: "elys1ve4a88p567vw3kq3uc20hp3pkqlrqjs2slgma3",
		Amount:  math.NewInt(71800000),
		Claimed: false,
	},
	{
		Address: "elys1tu7n2u2k72mrp36gh4akx9ljqdqc7zjd9ftu0k",
		Amount:  math.NewInt(71700000),
		Claimed: false,
	},
	{
		Address: "elys12pc593eqmug3srwks7mzu07rz9yswq88rxyetp",
		Amount:  math.NewInt(71600000),
		Claimed: false,
	},
	{
		Address: "elys13ynq6l6eafltrx7ymxfexyvlr0ye6e8g7mhtv6",
		Amount:  math.NewInt(71500000),
		Claimed: false,
	},
	{
		Address: "elys1hn58jnuqmklae0epv2t4a8kjhfr8nqec9wvc43",
		Amount:  math.NewInt(71500000),
		Claimed: false,
	},
	{
		Address: "elys1kf4q6ap73qzvj77xdv7jsqc45vcxd0gdtys5sj",
		Amount:  math.NewInt(71500000),
		Claimed: false,
	},
	{
		Address: "elys1lwnzyt697a4e7a7wwnz48tt9hzeqy0l0sy786d",
		Amount:  math.NewInt(71500000),
		Claimed: false,
	},
	{
		Address: "elys1z0dla9eazls2atfyqs3te8g3k9xtgsppueve07",
		Amount:  math.NewInt(71500000),
		Claimed: false,
	},
	{
		Address: "elys1dcd5yalcgsw4tfpqfdtq4k8j6e7y3g7f38j6vr",
		Amount:  math.NewInt(71400000),
		Claimed: false,
	},
	{
		Address: "elys1wtjans2pyrn4cjmpsryjk7763608zv7zmxwvuk",
		Amount:  math.NewInt(71300000),
		Claimed: false,
	},
	{
		Address: "elys16kx7z3s7axcsw8vyjdep09ffzujteulc84fwws",
		Amount:  math.NewInt(71100000),
		Claimed: false,
	},
	{
		Address: "elys1v2wc057r5kexfrnc2w6e0ag4w9vt2ycuvmpydk",
		Amount:  math.NewInt(71100000),
		Claimed: false,
	},
	{
		Address: "elys1jvu8jx4xzknn0m7enqpnzu734jke0nw9fztt74",
		Amount:  math.NewInt(71000000),
		Claimed: false,
	},
	{
		Address: "elys1c6d9292unx725h6epcxyr048q956qp3tcrqu4j",
		Amount:  math.NewInt(70800000),
		Claimed: false,
	},
	{
		Address: "elys1trr0e7kkl2tecwzdja6s9frkdyhnulatqfngsc",
		Amount:  math.NewInt(70500000),
		Claimed: false,
	},
	{
		Address: "elys1gemsr4dsrxvmpu2pgz0d3yz5em96e0qlmrclkd",
		Amount:  math.NewInt(70400000),
		Claimed: false,
	},
	{
		Address: "elys134ptlphnry4s8u778wspr3w8kwtgzwwavf06g2",
		Amount:  math.NewInt(70300000),
		Claimed: false,
	},
	{
		Address: "elys1hxsljdez9ak90jtq87z74v7v2ykcrlkthv09jx",
		Amount:  math.NewInt(70300000),
		Claimed: false,
	},
	{
		Address: "elys1pnm07d3960ny6c939dg88d6ld9gydumc7w8ljs",
		Amount:  math.NewInt(70300000),
		Claimed: false,
	},
	{
		Address: "elys1r8408d986r3gw6arlxgv2amc02sj5ulpt8t7xp",
		Amount:  math.NewInt(70300000),
		Claimed: false,
	},
	{
		Address: "elys1xgkaqnwfsw8cartfs8aaf9rjtg6wdvurmh7up9",
		Amount:  math.NewInt(70200000),
		Claimed: false,
	},
	{
		Address: "elys1vq9j5dr39xwsvee7n5wvxrnv4r93txplxzqh02",
		Amount:  math.NewInt(70000000),
		Claimed: false,
	},
	{
		Address: "elys1lvu7cmv2a5myeek2vqpq4926zsklka3ga6yyhz",
		Amount:  math.NewInt(69800000),
		Claimed: false,
	},
	{
		Address: "elys1zslgpud4gt0ys3dl35tzyrvw4mpk3scljgw6j7",
		Amount:  math.NewInt(69800000),
		Claimed: false,
	},
	{
		Address: "elys1vywd530rxwrn6x03z76nl2j9e3nw2r94wz2uxm",
		Amount:  math.NewInt(69700000),
		Claimed: false,
	},
	{
		Address: "elys1vam45gtu4pxvs3svv2sr3njyppzw8javknlq6n",
		Amount:  math.NewInt(69400000),
		Claimed: false,
	},
	{
		Address: "elys14n8r4d7zz6dpweqhpt7sqnaawthvce8nu0030q",
		Amount:  math.NewInt(69300000),
		Claimed: false,
	},
	{
		Address: "elys194mpe2m3gdg0mpfvxyj4tvttkhztjdl8gm43ld",
		Amount:  math.NewInt(69300000),
		Claimed: false,
	},
	{
		Address: "elys1xu0zheq9ejq8n604h6nmqfqmsp8c5wwew95h34",
		Amount:  math.NewInt(69200000),
		Claimed: false,
	},
	{
		Address: "elys1q50zgnxkzl94txqj0fygdtqgmmr5fdtnlpnks4",
		Amount:  math.NewInt(69000000),
		Claimed: false,
	},
	{
		Address: "elys1lhpexkj4fyr82836gguy492hwz2masscl0hhf9",
		Amount:  math.NewInt(68800000),
		Claimed: false,
	},
	{
		Address: "elys1lxxu2x9xtpc9277a5xhjk98ae0d6hk5p6x07jn",
		Amount:  math.NewInt(68800000),
		Claimed: false,
	},
	{
		Address: "elys1xmuqv3w35ym6kccscd9e49lfurhvzgs9pkmdyn",
		Amount:  math.NewInt(68800000),
		Claimed: false,
	},
	{
		Address: "elys104e983apdu47kd79mn8jha4z4470aauuhlfclp",
		Amount:  math.NewInt(68700000),
		Claimed: false,
	},
	{
		Address: "elys137nv58v5vesv4t0y0tqqqyksf6gqrum5adrh6f",
		Amount:  math.NewInt(68700000),
		Claimed: false,
	},
	{
		Address: "elys19udqt2z07ngetvf2xvq9ycuc0amdhxrm4c5r27",
		Amount:  math.NewInt(68700000),
		Claimed: false,
	},
	{
		Address: "elys1f4vpqwe0p0hctvu6nsnf59n04mx3uxfwdfs368",
		Amount:  math.NewInt(68700000),
		Claimed: false,
	},
	{
		Address: "elys1gkjqs06luwukpz5z6zketemcpzvzgu2wtw3vpq",
		Amount:  math.NewInt(68700000),
		Claimed: false,
	},
	{
		Address: "elys1jcs6ah3e97tfwxdkamsr84h7rkhwph7na77xfc",
		Amount:  math.NewInt(68700000),
		Claimed: false,
	},
	{
		Address: "elys1lr7zaympvay9892rwsztrhwvtpk3vzc8aafdhf",
		Amount:  math.NewInt(68700000),
		Claimed: false,
	},
	{
		Address: "elys1prry8nyvzuaw0hcj4au3lyd7gkwr4cw0njlx8p",
		Amount:  math.NewInt(68700000),
		Claimed: false,
	},
	{
		Address: "elys1t3ypmd38v3gtd46n3ahr788kcgxlulrs42v59f",
		Amount:  math.NewInt(68700000),
		Claimed: false,
	},
	{
		Address: "elys1tv9tlpwtj6rh5wg3eulspwzwlrn85cuenydca0",
		Amount:  math.NewInt(68700000),
		Claimed: false,
	},
	{
		Address: "elys1tw59eezlx5j8kr0grl3jryug44wafvmw0atr4p",
		Amount:  math.NewInt(68700000),
		Claimed: false,
	},
	{
		Address: "elys1yzd77f40r30kxftc3v9k3u7hguykgz3slwsksl",
		Amount:  math.NewInt(68600000),
		Claimed: false,
	},
	{
		Address: "elys15c4ny9ekq72t6pegv44x5zsugp3a5c5n2pg059",
		Amount:  math.NewInt(68500000),
		Claimed: false,
	},
	{
		Address: "elys1gzr6elf8cyslkd597ujapxgy3w300eqa0q2llf",
		Amount:  math.NewInt(68400000),
		Claimed: false,
	},
	{
		Address: "elys1lrhnrdr00y3tfctelge9vm5xe6dnzu9l9uquwk",
		Amount:  math.NewInt(68100000),
		Claimed: false,
	},
	{
		Address: "elys13mqs7ezr3en0g9snzuv0vu7q8kuv8xza8fehed",
		Amount:  math.NewInt(67900000),
		Claimed: false,
	},
	{
		Address: "elys1vpzlwskzcvy60lasrwa5xepztfp95vsdwfz92j",
		Amount:  math.NewInt(67900000),
		Claimed: false,
	},
	{
		Address: "elys1zfz0uguyn3t7nas94gza656rgeyzz9a7r23693",
		Amount:  math.NewInt(67900000),
		Claimed: false,
	},
	{
		Address: "elys19aq28vv8vqhn7mun27mtftesxf6v5vd7wmr4gq",
		Amount:  math.NewInt(67700000),
		Claimed: false,
	},
	{
		Address: "elys1zp9uaqctrewt9eklz4sxwkgp7jdvkuaytkrwrq",
		Amount:  math.NewInt(67700000),
		Claimed: false,
	},
	{
		Address: "elys18k2x9dg3uxjyxfspfz2mwlrjdwtse0hu2r5kvt",
		Amount:  math.NewInt(67600000),
		Claimed: false,
	},
	{
		Address: "elys1nefax734duw5je8g6y7vpp35756s9jsuuzt0e9",
		Amount:  math.NewInt(67500000),
		Claimed: false,
	},
	{
		Address: "elys1fqf4rzzd8wurwks7klgg588fm5gjt6tcvkr4ek",
		Amount:  math.NewInt(67300000),
		Claimed: false,
	},
	{
		Address: "elys15jn6ky9u48zyjxzksmu3pqv5ld9qclyfefn22h",
		Amount:  math.NewInt(67200000),
		Claimed: false,
	},
	{
		Address: "elys180gurwhyv0rlncqx706y9z0dgf7av0gjwhwgs0",
		Amount:  math.NewInt(67200000),
		Claimed: false,
	},
	{
		Address: "elys1fpm604q270j749q9nm8gqq7pdh6x5gdcnguudc",
		Amount:  math.NewInt(67200000),
		Claimed: false,
	},
	{
		Address: "elys1jz7av7cq45gh5hhrugtak7lkps2ga5v0m4wkk0",
		Amount:  math.NewInt(67200000),
		Claimed: false,
	},
	{
		Address: "elys17ad6tymv2f0hx5c5hqhzfmence44e52dvn5ta6",
		Amount:  math.NewInt(67099999),
		Claimed: false,
	},
	{
		Address: "elys1d82ufy9hmapfnk23jlta3qdrhzedq9khlrrflk",
		Amount:  math.NewInt(66900000),
		Claimed: false,
	},
	{
		Address: "elys18ml4ez6a7e0fpkx7qajynj0ha0getdwxz9q5q9",
		Amount:  math.NewInt(66700000),
		Claimed: false,
	},
	{
		Address: "elys1qyt6fqq9h3js7d0pmtam2zyq64mxhhzjfv5ssc",
		Amount:  math.NewInt(66700000),
		Claimed: false,
	},
	{
		Address: "elys1yyz66h9n2f6tay2u77c0djmf9g0kmvj6fvegr7",
		Amount:  math.NewInt(66599999),
		Claimed: false,
	},
	{
		Address: "elys1kmsteytxueay4s33567yqtk8ra8en84ae8lw5m",
		Amount:  math.NewInt(66500000),
		Claimed: false,
	},
	{
		Address: "elys1r7jm9zq6j26q2dtclv6za9ke9qp7e5vsa5wuqz",
		Amount:  math.NewInt(66400000),
		Claimed: false,
	},
	{
		Address: "elys1exnhc6nx0ulfkcxzqsj62fnfm3dutvtxwnamt5",
		Amount:  math.NewInt(66300000),
		Claimed: false,
	},
	{
		Address: "elys1f508pqu9rkeekkkfgh09zm7c4jla3nrxkq65l6",
		Amount:  math.NewInt(66300000),
		Claimed: false,
	},
	{
		Address: "elys1pmveqvhlz7e4yykanfak98xut2jaqyyx0kesry",
		Amount:  math.NewInt(66300000),
		Claimed: false,
	},
	{
		Address: "elys1t0hn9qt6cwz7er0la3eru4ye34z0awj6mg8klh",
		Amount:  math.NewInt(66300000),
		Claimed: false,
	},
	{
		Address: "elys1hmyahg20gpn8tah5cr86q34a08fusglwpz3scj",
		Amount:  math.NewInt(66200000),
		Claimed: false,
	},
	{
		Address: "elys1mk56848d2nhnapvvrr9a8fe8ny44gn2y0g2cx3",
		Amount:  math.NewInt(66099999),
		Claimed: false,
	},
	{
		Address: "elys1uq72qmr7gtsdj8k9cm6yfw5ry3g7ldj5uq34yk",
		Amount:  math.NewInt(66099999),
		Claimed: false,
	},
	{
		Address: "elys1xryzllx7cda8gf6gvmt6apddf4wdjkwysnpnx9",
		Amount:  math.NewInt(66099999),
		Claimed: false,
	},
	{
		Address: "elys14pv6m8sgmtx3z84rjef6dtg5rclnlzr880tln0",
		Amount:  math.NewInt(66000000),
		Claimed: false,
	},
	{
		Address: "elys1pthl6hc7r6ve2astpnpqhjl22xf9f36cxnh2aw",
		Amount:  math.NewInt(65900000),
		Claimed: false,
	},
	{
		Address: "elys1xezdal470gk9tt3ehqf66lw2jalejrkuquxtdk",
		Amount:  math.NewInt(65900000),
		Claimed: false,
	},
	{
		Address: "elys1yvvkkmzh8zs4z4v7gzn0ux8ag0me79ec3nvus9",
		Amount:  math.NewInt(65800000),
		Claimed: false,
	},
	{
		Address: "elys10nhsumc4ekpzp8cku9a686ygxwdh935ad787u8",
		Amount:  math.NewInt(65700000),
		Claimed: false,
	},
	{
		Address: "elys190les4zza48qv8fr0uwrkahhjtpgtjc4ard62v",
		Amount:  math.NewInt(65500000),
		Claimed: false,
	},
	{
		Address: "elys1esx7lec479qcl42dh33snqhtmz0nckd8krguxm",
		Amount:  math.NewInt(65400000),
		Claimed: false,
	},
	{
		Address: "elys16a7zl90xym79g5mf7l3m9va2zvrxtycc0ex8wc",
		Amount:  math.NewInt(65300000),
		Claimed: false,
	},
	{
		Address: "elys1m8rvmfrtwp90vgz573ktkqd2gmmw7w66g3ys80",
		Amount:  math.NewInt(65300000),
		Claimed: false,
	},
	{
		Address: "elys1rl5tpteu2g2gne2letw9zsctqcdutg34swkfnk",
		Amount:  math.NewInt(65300000),
		Claimed: false,
	},
	{
		Address: "elys1amlrs57zsfe8gngvp22jah5w6uhxx4urrcz22z",
		Amount:  math.NewInt(65200000),
		Claimed: false,
	},
	{
		Address: "elys1n43ucqmednjwzdsn6jd55xmnzzqe79a6ndtp38",
		Amount:  math.NewInt(65200000),
		Claimed: false,
	},
	{
		Address: "elys13xv0x0jxe4c3j2exkmada0p5ps45npknvwrr6e",
		Amount:  math.NewInt(65000000),
		Claimed: false,
	},
	{
		Address: "elys1zqtan0c2s82lt3m3f5jsrhl92064tly8csz960",
		Amount:  math.NewInt(65000000),
		Claimed: false,
	},
	{
		Address: "elys1zetezf50xyd8p4wc0525kcwakm7t0acacqgmuf",
		Amount:  math.NewInt(64800000),
		Claimed: false,
	},
	{
		Address: "elys15vehqrmgwwcdeut5q05c59e55w68ejzg7wa92h",
		Amount:  math.NewInt(64700000),
		Claimed: false,
	},
	{
		Address: "elys1szcfgdcv203t73ft5msj40qu7pep5gpkz6t8en",
		Amount:  math.NewInt(64700000),
		Claimed: false,
	},
	{
		Address: "elys1hdarmpvucjuy0hn39y27x5e5wr4kajz4c09qeu",
		Amount:  math.NewInt(64599999),
		Claimed: false,
	},
	{
		Address: "elys1v7ugxers4szc9xgnqz54atxvg5yfz7xsxchl22",
		Amount:  math.NewInt(64400000),
		Claimed: false,
	},
	{
		Address: "elys160gdptg53jj7enyxz3yf4py0slu2pfrv3dtj3c",
		Amount:  math.NewInt(64300000),
		Claimed: false,
	},
	{
		Address: "elys190vfx6t9dndkgmkhd228jprakapfg4q7tz2gtq",
		Amount:  math.NewInt(64300000),
		Claimed: false,
	},
	{
		Address: "elys1lsf9pq0w0zu4gdtye4ge8q6a0rjj373hxtrvl9",
		Amount:  math.NewInt(64200000),
		Claimed: false,
	},
	{
		Address: "elys1nuyqlen4f2g6yp84wszaudf2eetr92ms384juu",
		Amount:  math.NewInt(64000000),
		Claimed: false,
	},
	{
		Address: "elys1x050g2ancg6mshhhdheaaktcky9agwcpsz4xv3",
		Amount:  math.NewInt(64000000),
		Claimed: false,
	},
	{
		Address: "elys1ueq49f7whvnrhvlersermymf960a4h0nhlgkmm",
		Amount:  math.NewInt(63900000),
		Claimed: false,
	},
	{
		Address: "elys1k7szz79m4h4f03vdjd3jxzc2jc6796d9qux6rp",
		Amount:  math.NewInt(63800000),
		Claimed: false,
	},
	{
		Address: "elys1ecly4wtv8qgsf26n2eydr70zcsyfv893lfw9xc",
		Amount:  math.NewInt(63600000),
		Claimed: false,
	},
	{
		Address: "elys1yr9jazfqrl6xa22fcht73dgtgc4svz4cpakrfq",
		Amount:  math.NewInt(63600000),
		Claimed: false,
	},
	{
		Address: "elys1xrv89wtat29mkgr3y3hed7pn2a4uxp738mpa5p",
		Amount:  math.NewInt(63500000),
		Claimed: false,
	},
	{
		Address: "elys1hzpfqa23dfvfanhxl56xm544h2z3g0cq3eq49j",
		Amount:  math.NewInt(63200000),
		Claimed: false,
	},
	{
		Address: "elys1d35mtlpmc8z4pufs9pnk0wc8prqg4cw35grznx",
		Amount:  math.NewInt(63000000),
		Claimed: false,
	},
	{
		Address: "elys1k9qc4779thu7w2xc7qaec6j0jqfuxjc0wr3um5",
		Amount:  math.NewInt(62900000),
		Claimed: false,
	},
	{
		Address: "elys1g8p7jmq7gxmwvnazq2d40ygkqesu5denja2qve",
		Amount:  math.NewInt(62800000),
		Claimed: false,
	},
	{
		Address: "elys1zjtvs7na8mck7l2cyhdt8cnz2jzusxtdvr9qsy",
		Amount:  math.NewInt(62800000),
		Claimed: false,
	},
	{
		Address: "elys1tx4d0thz6djz5pdanfwjwe8643mhs0vm76fa4y",
		Amount:  math.NewInt(62700000),
		Claimed: false,
	},
	{
		Address: "elys1fl2yjvfvjn59g5amucudmvstufq08t34j9st95",
		Amount:  math.NewInt(62400000),
		Claimed: false,
	},
	{
		Address: "elys18sk8gxwpsadeqesht8vqhfanzh80xwwxrddn47",
		Amount:  math.NewInt(62300000),
		Claimed: false,
	},
	{
		Address: "elys1uwdcx6qzjhszqewwxe4ue429csrfxhpd3jzvtm",
		Amount:  math.NewInt(62300000),
		Claimed: false,
	},
	{
		Address: "elys1sjghvyzeyfw6gn37nc9l7wm4h236wjy8tavnr3",
		Amount:  math.NewInt(62200000),
		Claimed: false,
	},
	{
		Address: "elys1crkfmypdvlucckx8e67lxc0nr6wz5hmeeq8mpf",
		Amount:  math.NewInt(62100000),
		Claimed: false,
	},
	{
		Address: "elys1z0hwql0l0scmfne8805tyg62qf9p2utqg7dfeg",
		Amount:  math.NewInt(62100000),
		Claimed: false,
	},
	{
		Address: "elys10mk8954djw8xaxdqadjnu8rmttuwch3ay8rw7w",
		Amount:  math.NewInt(61900000),
		Claimed: false,
	},
	{
		Address: "elys15vdj0qxm9x7tzdmdnck70lp5u0rqpkn6z5l8pa",
		Amount:  math.NewInt(61900000),
		Claimed: false,
	},
	{
		Address: "elys139q6sqtxfhw2d9cdu2uk0drperfa0xyyad69ts",
		Amount:  math.NewInt(61800000),
		Claimed: false,
	},
	{
		Address: "elys1feh0v9vscsy0ej99gfnxzsy9f24mjar4w9x4yx",
		Amount:  math.NewInt(61800000),
		Claimed: false,
	},
	{
		Address: "elys1henzavr29335s898wgcmp2h7qld2s846ga527t",
		Amount:  math.NewInt(61800000),
		Claimed: false,
	},
	{
		Address: "elys18fstkxc5crfd725tvuhphgtrf3eru606fuzhvx",
		Amount:  math.NewInt(61700000),
		Claimed: false,
	},
	{
		Address: "elys1lmzgm4ulpp07wqhp86f4z7wjag9nlzskrvf34w",
		Amount:  math.NewInt(61600000),
		Claimed: false,
	},
	{
		Address: "elys1mklgzqqzmkjd0ma9d5ecssma9zftz60jhecdpn",
		Amount:  math.NewInt(61500000),
		Claimed: false,
	},
	{
		Address: "elys1gkme7wmetvpthp8c5nakxsngp4lku9p4qn3gw5",
		Amount:  math.NewInt(61400000),
		Claimed: false,
	},
	{
		Address: "elys10te4rzvxj3deym6cvvdndgkzrwl4ademgs6m6j",
		Amount:  math.NewInt(61300000),
		Claimed: false,
	},
	{
		Address: "elys15aqt4j8xt8wzkjf5z40n6ycpnh646e4885q83e",
		Amount:  math.NewInt(61300000),
		Claimed: false,
	},
	{
		Address: "elys1jltc6vrkr06x7vkfyfv0npvr905gkhuz5ucmla",
		Amount:  math.NewInt(61100000),
		Claimed: false,
	},
	{
		Address: "elys1vqxtaj48gg3xfkvp83re38ugrexpcy9nwf3wu0",
		Amount:  math.NewInt(60900000),
		Claimed: false,
	},
	{
		Address: "elys177r50vdy3jhe4tq536fzuur6u2mhpf76jers4k",
		Amount:  math.NewInt(60800000),
		Claimed: false,
	},
	{
		Address: "elys1e7xxgrwl6r5m0fe889u3sahknf03zxkhpk5tw6",
		Amount:  math.NewInt(60700000),
		Claimed: false,
	},
	{
		Address: "elys1nyw4lyw0qk6yccs9z58mhw8mu5tuxtu3rxpckg",
		Amount:  math.NewInt(60400000),
		Claimed: false,
	},
	{
		Address: "elys1k5zrehs2xl9yvrcrscwe2y283tav7klhf2xqgp",
		Amount:  math.NewInt(60300000),
		Claimed: false,
	},
	{
		Address: "elys1fd2wxqhwd7jtvsgkjpwjpkqjrh7j3mdwvpy7pm",
		Amount:  math.NewInt(60200000),
		Claimed: false,
	},
	{
		Address: "elys1ralltvn72rclqtd7szv7h8kckrw2s65t9f3fy4",
		Amount:  math.NewInt(60200000),
		Claimed: false,
	},
	{
		Address: "elys1ccldu966edm6gn225crdg2thhycxas5742wqzt",
		Amount:  math.NewInt(60100000),
		Claimed: false,
	},
	{
		Address: "elys1rs3v2p7ysye5wt2c55zeesafy0athu3sqxsg2v",
		Amount:  math.NewInt(60100000),
		Claimed: false,
	},
	{
		Address: "elys1032pwnhnzxdg3p9yjj7h3dqwj9z7600we34xas",
		Amount:  math.NewInt(60000000),
		Claimed: false,
	},
	{
		Address: "elys1tprrxsx0wwmf3vvmuvr6pn2nyvs6v49a3ytdhc",
		Amount:  math.NewInt(60000000),
		Claimed: false,
	},
	{
		Address: "elys1rdr5cx0v7gdhl6ep2n5vay75k56dru9nh3ew0w",
		Amount:  math.NewInt(59900000),
		Claimed: false,
	},
	{
		Address: "elys1d2lqlzkqsds0sr2u5pxd4ajka4aja7jkc4kvgy",
		Amount:  math.NewInt(59800000),
		Claimed: false,
	},
	{
		Address: "elys129tzkmzawxkn3lxdes5pzq436nx25lr0fj0lx3",
		Amount:  math.NewInt(59600000),
		Claimed: false,
	},
	{
		Address: "elys1smxq3wrxp7spkfgydh4dvn6025l6n7nkpwacrn",
		Amount:  math.NewInt(59600000),
		Claimed: false,
	},
	{
		Address: "elys1czfgvkdkmhq9x5j4ge5872qvumestq8zz8v06u",
		Amount:  math.NewInt(59300000),
		Claimed: false,
	},
	{
		Address: "elys1sq2jzlzxst39hqskmrkhy8cu75lncu9ry440l2",
		Amount:  math.NewInt(59300000),
		Claimed: false,
	},
	{
		Address: "elys13nhlg7cqspec5e0yv0n6nts3zcmfnc6h5qpww5",
		Amount:  math.NewInt(59000000),
		Claimed: false,
	},
	{
		Address: "elys1r79htf7p26w53vnerayv4ne9yfa6865sysk92u",
		Amount:  math.NewInt(59000000),
		Claimed: false,
	},
	{
		Address: "elys1fgevu5r55995zfkxdgw3jnxemqe46zzhp8ldtg",
		Amount:  math.NewInt(58900000),
		Claimed: false,
	},
	{
		Address: "elys1quuu0gn02any63mu3prk8s7rnumnjuhdtygm5l",
		Amount:  math.NewInt(58900000),
		Claimed: false,
	},
	{
		Address: "elys17dlkewww0h60zh3kylzpuphyyu33eyhkvc7h87",
		Amount:  math.NewInt(58700000),
		Claimed: false,
	},
	{
		Address: "elys1ytvrejvjmtvqcjr969shy9v8s6qgwcjfpr2lml",
		Amount:  math.NewInt(58700000),
		Claimed: false,
	},
	{
		Address: "elys19gvkw4srkva7qez0sjqsaa53qd9kpyf6kpqj0k",
		Amount:  math.NewInt(58400000),
		Claimed: false,
	},
	{
		Address: "elys14rmdxftvvqas547azl8tr2s9unynmjartjhtnj",
		Amount:  math.NewInt(58200000),
		Claimed: false,
	},
	{
		Address: "elys1e72v974wxev5mhlpf9a9eut56nvxwptr675q0f",
		Amount:  math.NewInt(58200000),
		Claimed: false,
	},
	{
		Address: "elys1gp4nsjt6pdcnx0qytnm8dlqkwve0yfua27aaup",
		Amount:  math.NewInt(58200000),
		Claimed: false,
	},
	{
		Address: "elys1k24uexr8ndmjx5tq5dy3xs48uxplcx759lh5wn",
		Amount:  math.NewInt(58200000),
		Claimed: false,
	},
	{
		Address: "elys1vp3v7ncrt0anu6tnwazgxwgz32u2fel7556q46",
		Amount:  math.NewInt(58200000),
		Claimed: false,
	},
	{
		Address: "elys1m3sgh6zrafma5vqjnx9snepfdyr78tvzzpqkyr",
		Amount:  math.NewInt(58100000),
		Claimed: false,
	},
	{
		Address: "elys1wunzut4f7mg22j47es2y989n5nzm8h8w8lkewy",
		Amount:  math.NewInt(58000000),
		Claimed: false,
	},
	{
		Address: "elys12djqcjectg3zqe08ajryd4ahddayz8e2djw4gx",
		Amount:  math.NewInt(57900000),
		Claimed: false,
	},
	{
		Address: "elys1gjx46yqqwkvwh3rc8gv3j80hg496ejs9xzf4uu",
		Amount:  math.NewInt(57900000),
		Claimed: false,
	},
	{
		Address: "elys1umj73q5kpw7ag4qc5q7y5k94l9lfcp3pqzzrqh",
		Amount:  math.NewInt(57800000),
		Claimed: false,
	},
	{
		Address: "elys1vcwqz59qf0dh7hgwy64dm2eyyqwp49d0h3kje2",
		Amount:  math.NewInt(57700000),
		Claimed: false,
	},
	{
		Address: "elys10trtua6hkekzsekpa05ykyr6eveej5lqhtlzqk",
		Amount:  math.NewInt(57600000),
		Claimed: false,
	},
	{
		Address: "elys1shwck77lgnle66sg9satc0r0vwu5avs6d6kllx",
		Amount:  math.NewInt(57600000),
		Claimed: false,
	},
	{
		Address: "elys1sdhll32xy5yrcn96tr5yvu0uyy5qr2fja86kvs",
		Amount:  math.NewInt(57500000),
		Claimed: false,
	},
	{
		Address: "elys1c97hke0mhexhc60us7c7xdw2pktspyj6vqt2pt",
		Amount:  math.NewInt(57400000),
		Claimed: false,
	},
	{
		Address: "elys1usgh9pthhpgjaulw3sh2gwppsnnypgkju9n43s",
		Amount:  math.NewInt(57400000),
		Claimed: false,
	},
	{
		Address: "elys1h8pd589q68ltjhyl7mycq3xq77ynsxlqe8mad9",
		Amount:  math.NewInt(57300000),
		Claimed: false,
	},
	{
		Address: "elys1qdywps48z4erm8tyckwhmzdwtxd45rke6jfh4y",
		Amount:  math.NewInt(57200000),
		Claimed: false,
	},
	{
		Address: "elys1z87p546v29u7hydhgm2ndkalmjxnztpw53m6xw",
		Amount:  math.NewInt(57000000),
		Claimed: false,
	},
	{
		Address: "elys164ycpphrs8cwhdm44404hfs8h24v7evn0tynqs",
		Amount:  math.NewInt(56700000),
		Claimed: false,
	},
	{
		Address: "elys1c80v3xteftfnpgen0uhsjr4zk06f87nq3t3hks",
		Amount:  math.NewInt(56700000),
		Claimed: false,
	},
	{
		Address: "elys1hexpssha27sn5gwma4qg2uywa23y4sp6kvxe7c",
		Amount:  math.NewInt(56700000),
		Claimed: false,
	},
	{
		Address: "elys1q7aj38ln6gs6kn83kvu304qcr8kpswu9hzkjer",
		Amount:  math.NewInt(56600000),
		Claimed: false,
	},
	{
		Address: "elys1dcgqgkj4hmf30u0z9v6mnsvres67cxww7gvasw",
		Amount:  math.NewInt(56300000),
		Claimed: false,
	},
	{
		Address: "elys1qw0a4v4fexfsdt93kzlestd0lr5arwpeksnwl8",
		Amount:  math.NewInt(56300000),
		Claimed: false,
	},
	{
		Address: "elys1vm47td49lslp76yzhyel3wsvzvp5yea4yyvew8",
		Amount:  math.NewInt(56300000),
		Claimed: false,
	},
	{
		Address: "elys1303986x7yqr7fux0w364hs9z30zt8ydd8t9krh",
		Amount:  math.NewInt(56200000),
		Claimed: false,
	},
	{
		Address: "elys169fugs8dxxp6smsf3n47m4eqr5z95chx8ljq72",
		Amount:  math.NewInt(56200000),
		Claimed: false,
	},
	{
		Address: "elys18vsxav4z2fa9ykdj29h3c6xyauu4xl5x8ncqfq",
		Amount:  math.NewInt(56200000),
		Claimed: false,
	},
	{
		Address: "elys1mlyccp906auy36ytqd6pkjta50htcjcz5c4q4g",
		Amount:  math.NewInt(56200000),
		Claimed: false,
	},
	{
		Address: "elys1tk4hcnnx3xtcc5k50juap4699548x8actlt583",
		Amount:  math.NewInt(56200000),
		Claimed: false,
	},
	{
		Address: "elys1jxe3vn6e6ahc65ak9p4t3sunuvy6vqhjeqm6am",
		Amount:  math.NewInt(56100000),
		Claimed: false,
	},
	{
		Address: "elys1tkr73y8eyevv7du2ennpej8nsh2hlvhqwddqls",
		Amount:  math.NewInt(56100000),
		Claimed: false,
	},
	{
		Address: "elys14phy6ckg6pv2pvpw2e53pt7n2ja4rvf9egwmk2",
		Amount:  math.NewInt(56000000),
		Claimed: false,
	},
	{
		Address: "elys1ezdlgqfq6u88e57y3rjatdccw44cmtwk9t258j",
		Amount:  math.NewInt(55900000),
		Claimed: false,
	},
	{
		Address: "elys1vdneccvnfmlhxj85fgxkrvupuftyzem48avmgq",
		Amount:  math.NewInt(55800000),
		Claimed: false,
	},
	{
		Address: "elys1ywl67ewjhn9s2wdx06amfuv838epkgjhrrdjq6",
		Amount:  math.NewInt(55800000),
		Claimed: false,
	},
	{
		Address: "elys1c4dmm2xwtfyz3yqhg3ya3supgm0zsf7px6s6f7",
		Amount:  math.NewInt(55600000),
		Claimed: false,
	},
	{
		Address: "elys1gu6a05ttjshwn2yhd67uw2vdn3j4csww7e0hym",
		Amount:  math.NewInt(55600000),
		Claimed: false,
	},
	{
		Address: "elys1j9snzzd732asn8wuv5llxvq3dfe3lapqjgmuq3",
		Amount:  math.NewInt(55600000),
		Claimed: false,
	},
	{
		Address: "elys182cgd09j3g9a9djhu2dvdd7a4nwq3y6mselg5r",
		Amount:  math.NewInt(55500000),
		Claimed: false,
	},
	{
		Address: "elys1s2utctaw56k2vndrd6aw0wmkne47ymda8g28ug",
		Amount:  math.NewInt(55500000),
		Claimed: false,
	},
	{
		Address: "elys1pqvwygcc9d022yef2xlwkdslh8zkxc5lep7d9a",
		Amount:  math.NewInt(55000000),
		Claimed: false,
	},
	{
		Address: "elys1v28ce475qp60e6ks092huay2fpsrz385r0kfnn",
		Amount:  math.NewInt(55000000),
		Claimed: false,
	},
	{
		Address: "elys13z62f56zd4q3et9xkuc8f3xhnk8324z7qen4lf",
		Amount:  math.NewInt(54800000),
		Claimed: false,
	},
	{
		Address: "elys1lqfn380a40evf8vt3f2e8y0vcdsmuvjal3yq40",
		Amount:  math.NewInt(54800000),
		Claimed: false,
	},
	{
		Address: "elys14wt9d6xj7n7gk9dwmd8lep2tnp4jjk7tnjwlwk",
		Amount:  math.NewInt(54600000),
		Claimed: false,
	},
	{
		Address: "elys1z03wkekwnafvqzccwq2t47y9yajyh3dl9hle23",
		Amount:  math.NewInt(54600000),
		Claimed: false,
	},
	{
		Address: "elys1ruwraamthss0nsvgye7pwnl6jarxwk0v86skc3",
		Amount:  math.NewInt(54500000),
		Claimed: false,
	},
	{
		Address: "elys137l9vr6dvfn3pqq8yjxj04njaneldruawferh2",
		Amount:  math.NewInt(54400000),
		Claimed: false,
	},
	{
		Address: "elys1ue03088ptzf8tze4z7saxzjp223uxwmx0lshfy",
		Amount:  math.NewInt(54200000),
		Claimed: false,
	},
	{
		Address: "elys1wvaz2qvxhl7jpgfm2q38257z8sv7makf2avcwm",
		Amount:  math.NewInt(54100000),
		Claimed: false,
	},
	{
		Address: "elys12ruvl0j8nwvvjc58znlvvsaqkefmplda39q70e",
		Amount:  math.NewInt(54000000),
		Claimed: false,
	},
	{
		Address: "elys1ckeqcdqxls7cr5azkj62ywaprvh8gr04h9s5sm",
		Amount:  math.NewInt(54000000),
		Claimed: false,
	},
	{
		Address: "elys1wfq4thzedvaflrv996zrhzs03ed5mr4xe7vlz7",
		Amount:  math.NewInt(53900000),
		Claimed: false,
	},
	{
		Address: "elys1ydt0jchl3sed37xuuprp4xhczx6j4umd0hw6xs",
		Amount:  math.NewInt(53900000),
		Claimed: false,
	},
	{
		Address: "elys1c8x9trxfy6hhx8qsll2q6zzdsfp32f3k9py9au",
		Amount:  math.NewInt(53800000),
		Claimed: false,
	},
	{
		Address: "elys1pujxe7ggahkppufvg880seefhqu94zex4ahzy4",
		Amount:  math.NewInt(53800000),
		Claimed: false,
	},
	{
		Address: "elys1s2csffg7m49tkwv3e9xxtvw0hzp5l5fcjn9yy3",
		Amount:  math.NewInt(53800000),
		Claimed: false,
	},
	{
		Address: "elys1l9krrqwe0hheck5gvjy29klkk5vfpt3y59tzk3",
		Amount:  math.NewInt(53700000),
		Claimed: false,
	},
	{
		Address: "elys10p23qgwnsw5wlls4vlge4elzgxsdwncdmvv55q",
		Amount:  math.NewInt(53600000),
		Claimed: false,
	},
	{
		Address: "elys130pqjgslkldgpvckncuf345newxwu0ga48d6n5",
		Amount:  math.NewInt(53400000),
		Claimed: false,
	},
	{
		Address: "elys1r2ncj8wth5f4huhk5nj09hg92czdh7vyefl6qx",
		Amount:  math.NewInt(53400000),
		Claimed: false,
	},
	{
		Address: "elys16xqf05lggllk363x7gppll5x8jh7yz89ymuf5j",
		Amount:  math.NewInt(53300000),
		Claimed: false,
	},
	{
		Address: "elys1c2vmxp4y5hf5h2scvxqpmrtdu7ned82uprrtwg",
		Amount:  math.NewInt(53200000),
		Claimed: false,
	},
	{
		Address: "elys1h5n0g98fqsf9xluzw9k23nhqurmjyx3tsprkwz",
		Amount:  math.NewInt(53200000),
		Claimed: false,
	},
	{
		Address: "elys1vduy6g4f2qckgcn30n9atrhqy9zuq5939g5g3f",
		Amount:  math.NewInt(53200000),
		Claimed: false,
	},
	{
		Address: "elys1l7knujvc4prnvee3scas7vnvnv2mq7defvejuh",
		Amount:  math.NewInt(53100000),
		Claimed: false,
	},
	{
		Address: "elys1puzevxeu79366dcg6ffaffha9uf79h0fg53un3",
		Amount:  math.NewInt(53000000),
		Claimed: false,
	},
	{
		Address: "elys1sjed9gjs3lacn4jvz3ur3q097z3y4hfhsd3l57",
		Amount:  math.NewInt(52900000),
		Claimed: false,
	},
	{
		Address: "elys12ht53cjjgdtlxgcyuexc7jqd2lkexufphlcw2l",
		Amount:  math.NewInt(52700000),
		Claimed: false,
	},
	{
		Address: "elys15f7vwa5vr2gamyugrh80jggpmyn6pyx6mxzg8g",
		Amount:  math.NewInt(52600000),
		Claimed: false,
	},
	{
		Address: "elys1jltnut65enecswfhnvuemywm4vdvtahw564zz9",
		Amount:  math.NewInt(52600000),
		Claimed: false,
	},
	{
		Address: "elys1yspal5ahp8nqyfut6rr5fggk8pltwhfd0md0q3",
		Amount:  math.NewInt(52600000),
		Claimed: false,
	},
	{
		Address: "elys1942ls89gfu5axflcakvwm799y6vue327tg557m",
		Amount:  math.NewInt(52300000),
		Claimed: false,
	},
	{
		Address: "elys1kectm60a0cmrh3cl5ld33vt59pu37nkjqdw0gl",
		Amount:  math.NewInt(52300000),
		Claimed: false,
	},
	{
		Address: "elys1tc5pt52kjs6d9hk2xgf3c9a32tsnkdf4xz48fz",
		Amount:  math.NewInt(52300000),
		Claimed: false,
	},
	{
		Address: "elys1lkmzlxw5x3u9tcmyp7u8qcmvz9m05r7g06clll",
		Amount:  math.NewInt(52200000),
		Claimed: false,
	},
	{
		Address: "elys1hr84q67algxx8xvaegh0pacddv0je6x0g3tqde",
		Amount:  math.NewInt(52100000),
		Claimed: false,
	},
	{
		Address: "elys1nqawz9nf3jh7u9uwt5c4d8ncupyfth438yy7q0",
		Amount:  math.NewInt(52100000),
		Claimed: false,
	},
	{
		Address: "elys1pcz5vr8mcx6w8nxgll8s0m77vtnqdhc3fe06e4",
		Amount:  math.NewInt(51900000),
		Claimed: false,
	},
	{
		Address: "elys1wd32n7k6qdqssgujajwxehgudd2jf864d45vfk",
		Amount:  math.NewInt(51900000),
		Claimed: false,
	},
	{
		Address: "elys1qj2y386lgz4x2g8lmrjxe909natp2atukxskdk",
		Amount:  math.NewInt(51800000),
		Claimed: false,
	},
	{
		Address: "elys18qtemd008rnrpzuszex4angzqq8753fhnnpqnz",
		Amount:  math.NewInt(51700000),
		Claimed: false,
	},
	{
		Address: "elys1nhjygea5fvf2lkd9u50sdpp3elc0hdlchg8s66",
		Amount:  math.NewInt(51700000),
		Claimed: false,
	},
	{
		Address: "elys178t3y3387nu5eumyk0k50mqurx9c947888jj2l",
		Amount:  math.NewInt(51600000),
		Claimed: false,
	},
	{
		Address: "elys1lvgxah9ppcvcshngk4kcquwwx59ms7hclszpmn",
		Amount:  math.NewInt(51600000),
		Claimed: false,
	},
	{
		Address: "elys1434ashl4wqx9pm5vtg2v929utp5ndca878tch3",
		Amount:  math.NewInt(51500000),
		Claimed: false,
	},
	{
		Address: "elys14e04xex0yg635r46u36l2l320t980xdhn5dfat",
		Amount:  math.NewInt(51500000),
		Claimed: false,
	},
	{
		Address: "elys1l2yqyp9zqa6q2gw5sxxhz86ugkttxmmzapxkku",
		Amount:  math.NewInt(51500000),
		Claimed: false,
	},
	{
		Address: "elys1tw7wfucn2zfceecpy3fuf3tq7naw8dvssevylg",
		Amount:  math.NewInt(51500000),
		Claimed: false,
	},
	{
		Address: "elys1ujea0zf3y3mvfgk5m0nnj9r94kqhwmc658dsql",
		Amount:  math.NewInt(51300000),
		Claimed: false,
	},
	{
		Address: "elys1ylcf89f4gklt96j0a4jwsu7gfrlh8j9e7axmqm",
		Amount:  math.NewInt(51300000),
		Claimed: false,
	},
	{
		Address: "elys14c6mgrcc9rnfzxcrptyecj003l0cjv0h0z3rtd",
		Amount:  math.NewInt(51200000),
		Claimed: false,
	},
	{
		Address: "elys1crgy3sja7gu2t37wnsxl4p3j95cnxg3jg0yt54",
		Amount:  math.NewInt(51200000),
		Claimed: false,
	},
	{
		Address: "elys13x6eyr23zk54j5sexhastzwu2hg27gfy848wdu",
		Amount:  math.NewInt(51000000),
		Claimed: false,
	},
	{
		Address: "elys1737jyr99x4mezjms4zffajp42hclckypqu2uca",
		Amount:  math.NewInt(50800000),
		Claimed: false,
	},
	{
		Address: "elys1wq6njy3944tuhzekhpyswpwn8cs2w5syw7zq7q",
		Amount:  math.NewInt(50800000),
		Claimed: false,
	},
	{
		Address: "elys1mxd5t3qtte3t43rzsspu25834rz5c6wyvvx25m",
		Amount:  math.NewInt(50600000),
		Claimed: false,
	},
	{
		Address: "elys1f3v68c32nw8c4xhqt59cugk5kwsv2wmkqcly6c",
		Amount:  math.NewInt(50400000),
		Claimed: false,
	},
	{
		Address: "elys1qpsdklr4pw5r4ykz3mfrjreeatdx3v7q5ftvpc",
		Amount:  math.NewInt(50400000),
		Claimed: false,
	},
	{
		Address: "elys15s4h0ezv3eqv2k78e5cy2x3pt3w36jwjqkztpq",
		Amount:  math.NewInt(50300000),
		Claimed: false,
	},
	{
		Address: "elys1d4c4xk399fa9tvztsa82hqk2wh0n3ygvkp2hja",
		Amount:  math.NewInt(50300000),
		Claimed: false,
	},
	{
		Address: "elys1mxrunu28tngaqydt4d5d4vr7altkpc7j2sn6cd",
		Amount:  math.NewInt(50300000),
		Claimed: false,
	},
	{
		Address: "elys1tk9wvp4glsxz0n2plrxa7qhukpe2e22jllag22",
		Amount:  math.NewInt(50300000),
		Claimed: false,
	},
	{
		Address: "elys1qdvh69au2953d3gr0ujudz6p6rhyyyxejl5qr8",
		Amount:  math.NewInt(50200000),
		Claimed: false,
	},
	{
		Address: "elys1s7g66zqszetffcypu5n37ckx5aly4tk0fvnvvw",
		Amount:  math.NewInt(50200000),
		Claimed: false,
	},
	{
		Address: "elys15927y2rj2fhevpawxk3aatask53rj5c9e7ds69",
		Amount:  math.NewInt(50100000),
		Claimed: false,
	},
	{
		Address: "elys17vuy9sgwr39td5z27ffjswpyuk4qs4p0yazes4",
		Amount:  math.NewInt(50100000),
		Claimed: false,
	},
	{
		Address: "elys1d0g3m77g2gm4rm7qw2jtsfqthergc73yn6ze4f",
		Amount:  math.NewInt(50100000),
		Claimed: false,
	},
	{
		Address: "elys1f53zhyfds4yjv9dlcw32ahck3xjy2j65vkrd8d",
		Amount:  math.NewInt(50100000),
		Claimed: false,
	},
	{
		Address: "elys1ncytyhx2pn0khye9cxavd9lsrvzcs9nz9z0pqy",
		Amount:  math.NewInt(50000000),
		Claimed: false,
	},
	{
		Address: "elys1y03src9ej2mtnxrmw5w2d7laljhf3mljezxruf",
		Amount:  math.NewInt(50000000),
		Claimed: false,
	},
	{
		Address: "elys1635gyt647tch2yjvkfjvwdxj97e7vceygd4uqp",
		Amount:  math.NewInt(49900000),
		Claimed: false,
	},
	{
		Address: "elys1u62ksu6vl4nke79hzvgv78jm4vav783zphfnpt",
		Amount:  math.NewInt(49900000),
		Claimed: false,
	},
	{
		Address: "elys178nrcg9r9etnfxhcn86s224csy52zpcqkrwk83",
		Amount:  math.NewInt(49700000),
		Claimed: false,
	},
	{
		Address: "elys1xt8yng6mke3l8880m85awf66dquzx25t24vgw9",
		Amount:  math.NewInt(49700000),
		Claimed: false,
	},
	{
		Address: "elys15rfdf2u9z9mng50npjtqxkuglqyzvax7vp8aar",
		Amount:  math.NewInt(49500000),
		Claimed: false,
	},
	{
		Address: "elys19g4r3x8rycchd7k8tarq34367kmsn7fw7mpjq9",
		Amount:  math.NewInt(49400000),
		Claimed: false,
	},
	{
		Address: "elys1tkwlz65yp0c9l2zmcdawtepgvpmvgkaa4ydpdd",
		Amount:  math.NewInt(49400000),
		Claimed: false,
	},
	{
		Address: "elys16zlzmce8klra8je5qu24q8pw7tgplnt66qy2mk",
		Amount:  math.NewInt(49300000),
		Claimed: false,
	},
	{
		Address: "elys19lqhe9rr2zgjheqttg9kg33herf2d89a6d0cyy",
		Amount:  math.NewInt(49300000),
		Claimed: false,
	},
	{
		Address: "elys1lz6nesp98c4l5ms9k3edn505r9fldcdppgq82q",
		Amount:  math.NewInt(49300000),
		Claimed: false,
	},
	{
		Address: "elys1a9ql8jkw4wuj2jqdgnnjaza0jmfk4a7uj5psx6",
		Amount:  math.NewInt(49200000),
		Claimed: false,
	},
	{
		Address: "elys1htghsz7ew56u2tvhpwszpq72y6xlp36mem6rru",
		Amount:  math.NewInt(49200000),
		Claimed: false,
	},
	{
		Address: "elys14kr56vks0hnphza2ajfcktggjvssw45ylxee5u",
		Amount:  math.NewInt(49100000),
		Claimed: false,
	},
	{
		Address: "elys1n0pd4qhxl8u55qc0uanluu8ar2q8fukcjrr82d",
		Amount:  math.NewInt(49100000),
		Claimed: false,
	},
	{
		Address: "elys192sualvmc8jy5yzs79vey9kx3n6yawe0ag5u0e",
		Amount:  math.NewInt(48900000),
		Claimed: false,
	},
	{
		Address: "elys1jvnpaayv96qfzgls36yggqnvmv5ssxtf97km0j",
		Amount:  math.NewInt(48700000),
		Claimed: false,
	},
	{
		Address: "elys1kaj7x6jjmuzku0t5ze5g8dsgh8n9n8wj4snhqj",
		Amount:  math.NewInt(48700000),
		Claimed: false,
	},
	{
		Address: "elys1agg40pllc6d6qnd57eu50rgmsqyc4vmswnkes6",
		Amount:  math.NewInt(48500000),
		Claimed: false,
	},
	{
		Address: "elys137kzfp6kjeamyswhjwkj3jux806thzdg4jen2u",
		Amount:  math.NewInt(48400000),
		Claimed: false,
	},
	{
		Address: "elys19ly7nelc834rdh96sy0s3k7q455gn9sn3lfrlf",
		Amount:  math.NewInt(48400000),
		Claimed: false,
	},
	{
		Address: "elys1w6wp86lce7um7ddc289askca0yydnx9fgexyxx",
		Amount:  math.NewInt(48200000),
		Claimed: false,
	},
	{
		Address: "elys12e03drfhq5rz8wtjaag83pl4pxy7n9tg8rh32x",
		Amount:  math.NewInt(48100000),
		Claimed: false,
	},
	{
		Address: "elys1nxj7wfrvjqmwqp9m9t7pft7u3fekflcdvy3gz6",
		Amount:  math.NewInt(48100000),
		Claimed: false,
	},
	{
		Address: "elys1pwdj562q24sk0ktqpaeek3n7eknrxv6r4k7286",
		Amount:  math.NewInt(48100000),
		Claimed: false,
	},
	{
		Address: "elys1qsvafpdelsjrxhvrlgcpfz0d368kckhn79dcq7",
		Amount:  math.NewInt(48000000),
		Claimed: false,
	},
	{
		Address: "elys132gxdzvfzj4nj8red94vv37adrtggu8ae02vt9",
		Amount:  math.NewInt(47800000),
		Claimed: false,
	},
	{
		Address: "elys1rgnd3wftjyzm4j0jgkxlcvfhs4u6qafwtcwazq",
		Amount:  math.NewInt(47700000),
		Claimed: false,
	},
	{
		Address: "elys1e04xqhcstqxvesvys50h7k90tgfvwvpfa0a2q4",
		Amount:  math.NewInt(47500000),
		Claimed: false,
	},
	{
		Address: "elys1gms663jp50ufvwgckppf2mdu8g3x6d6lkzde3n",
		Amount:  math.NewInt(47400000),
		Claimed: false,
	},
	{
		Address: "elys1mq8tg937ufhzxj4tgps4y0vjg4l8v3cmellqlw",
		Amount:  math.NewInt(47400000),
		Claimed: false,
	},
	{
		Address: "elys1uymaerenzyz47l9aryjlnll5ecn0vq552pq2rv",
		Amount:  math.NewInt(47400000),
		Claimed: false,
	},
	{
		Address: "elys1tcfez7gk36vg7jpq474796602ljtkq5wvgy6qy",
		Amount:  math.NewInt(47200000),
		Claimed: false,
	},
	{
		Address: "elys1xqs5fhw6xecqvcs29j2w4y0p4438qk5j32rg2w",
		Amount:  math.NewInt(47200000),
		Claimed: false,
	},
	{
		Address: "elys1mklsh50f22fq8h2cfn2crmu5kqsuw2nxn7pxs0",
		Amount:  math.NewInt(47100000),
		Claimed: false,
	},
	{
		Address: "elys13hxjm03qdxtqd70r3r4sc77sz6xzgmc8q9n822",
		Amount:  math.NewInt(47000000),
		Claimed: false,
	},
	{
		Address: "elys14k47vh8fn54xu2j59pfa2tuykn0jmqsd4pyax2",
		Amount:  math.NewInt(47000000),
		Claimed: false,
	},
	{
		Address: "elys1pfxqjrlcfvyexpzklxextstd7z332era22tt2k",
		Amount:  math.NewInt(47000000),
		Claimed: false,
	},
	{
		Address: "elys1gpsfclgpz58fun3jtqr2gpq3r5cuhq3rxtnpzf",
		Amount:  math.NewInt(46600000),
		Claimed: false,
	},
	{
		Address: "elys1lv0lq9aw5dcfj42u9yl3vvhsy48j0v9kx57sjc",
		Amount:  math.NewInt(46600000),
		Claimed: false,
	},
	{
		Address: "elys1zsxtdtetr5ahzw25qyeugs5aggwfqqmh6tspck",
		Amount:  math.NewInt(46400000),
		Claimed: false,
	},
	{
		Address: "elys10me5lgjppdazsu8u87kt2j32zju0j5v5q2qjuk",
		Amount:  math.NewInt(46300000),
		Claimed: false,
	},
	{
		Address: "elys1pwegtads76mknhy4kxe5u7gauckw6q7kvdq2ke",
		Amount:  math.NewInt(46300000),
		Claimed: false,
	},
	{
		Address: "elys16q5f4etzsc6lhm7xsncc7ne6rxwlaeg3jjljyh",
		Amount:  math.NewInt(46200000),
		Claimed: false,
	},
	{
		Address: "elys1thg5dqlapylvtvgxard2rpgefrvdsukal9v4gy",
		Amount:  math.NewInt(46200000),
		Claimed: false,
	},
	{
		Address: "elys1ff5gxykrqasuj3k689kk7ley7hz552e56cmm0s",
		Amount:  math.NewInt(45800000),
		Claimed: false,
	},
	{
		Address: "elys1x08lvfymh34zkhu0en45we73gcmwhee9qcn7l6",
		Amount:  math.NewInt(45800000),
		Claimed: false,
	},
	{
		Address: "elys1x4t29yjxudh4wt3lfx0ed4fcu94360fu5nysyf",
		Amount:  math.NewInt(45800000),
		Claimed: false,
	},
	{
		Address: "elys10myu4h4774n7clc7647rgwpqcqe6kmw5vw55e7",
		Amount:  math.NewInt(45700000),
		Claimed: false,
	},
	{
		Address: "elys1cs0ysv4l5rk6h3tayrza4x9lphk8dpjvnzgzl5",
		Amount:  math.NewInt(45700000),
		Claimed: false,
	},
	{
		Address: "elys1g4r8tvxtd3xzv5kjq7yez6e2ej7phldyu6ce3m",
		Amount:  math.NewInt(45700000),
		Claimed: false,
	},
	{
		Address: "elys1uzhhp4jjhhcq2vr4zn022ukf5tn9fekk7fxmrp",
		Amount:  math.NewInt(45700000),
		Claimed: false,
	},
	{
		Address: "elys13ysztw5lt2drxlrkyw9f3jx7u30c63avlal5ku",
		Amount:  math.NewInt(45500000),
		Claimed: false,
	},
	{
		Address: "elys1a3qz9c9mmw76pvzka8p2pj32fncvkl0s6xmtka",
		Amount:  math.NewInt(45500000),
		Claimed: false,
	},
	{
		Address: "elys1vx48md64dnv6v3j3uum2fae49cgxctswnhx0ee",
		Amount:  math.NewInt(45500000),
		Claimed: false,
	},
	{
		Address: "elys1wv5a4quvyg7rugg99wj9l2n8ardg6ccys8tsdp",
		Amount:  math.NewInt(45500000),
		Claimed: false,
	},
	{
		Address: "elys13uergunccum7hz7hwrtun4dnqz4pe8zg3v84t2",
		Amount:  math.NewInt(45300000),
		Claimed: false,
	},
	{
		Address: "elys1ewt6wnzspq8xqu0ytpn30484jd8hv5h52c25nk",
		Amount:  math.NewInt(45300000),
		Claimed: false,
	},
	{
		Address: "elys1g4ljd7h8ft03gf3jvvyfv6gfmkg7tfj49x4man",
		Amount:  math.NewInt(45300000),
		Claimed: false,
	},
	{
		Address: "elys1svf58wkyn9jv0clqwpmevydqyx0ch4lyhprm28",
		Amount:  math.NewInt(45300000),
		Claimed: false,
	},
	{
		Address: "elys1wujzpzg2yjlvcuy6hvlq5tkdchtmkrjzskk922",
		Amount:  math.NewInt(45300000),
		Claimed: false,
	},
	{
		Address: "elys1kk53s2ylp7ge65stvjvq2s63q0s72890apjvmn",
		Amount:  math.NewInt(45200000),
		Claimed: false,
	},
	{
		Address: "elys1mmfp9cwmh6v38qc3ny2sppwc6yhzt0zy73ctue",
		Amount:  math.NewInt(45200000),
		Claimed: false,
	},
	{
		Address: "elys1pdwdxu4wwvqatdjpp5e7yx0wy37tjherm7mjwp",
		Amount:  math.NewInt(45200000),
		Claimed: false,
	},
	{
		Address: "elys10tdeh07ep4e7tm92ez87ad842t05wzamrdjr55",
		Amount:  math.NewInt(45100000),
		Claimed: false,
	},
	{
		Address: "elys1axnryg6wt6pkqqcqrtsqrz5a9djx0sgkfq0prn",
		Amount:  math.NewInt(45100000),
		Claimed: false,
	},
	{
		Address: "elys1kez0dypydmzc8a09sptarw95mjepfk74x3uaxd",
		Amount:  math.NewInt(45100000),
		Claimed: false,
	},
	{
		Address: "elys1ekxyan4a8hldlzyfrpt2h3j4ajt9lcxk2ufhj2",
		Amount:  math.NewInt(45000000),
		Claimed: false,
	},
	{
		Address: "elys1zv8qvgwxs5r4cfqqhtw2wjv53u4vccfetvpyq8",
		Amount:  math.NewInt(44900000),
		Claimed: false,
	},
	{
		Address: "elys1xmvp95nz3lmfumr2h5shmxd7jzyd60eszhqxmt",
		Amount:  math.NewInt(44800000),
		Claimed: false,
	},
	{
		Address: "elys17mdmyzwm3ja50fc446txad5vxcmvtc7em2fr55",
		Amount:  math.NewInt(44700000),
		Claimed: false,
	},
	{
		Address: "elys1sgalu8meqw89xvtndp50qhd3p0r4hsas2rmy8p",
		Amount:  math.NewInt(44700000),
		Claimed: false,
	},
	{
		Address: "elys1yrnh6z3r3ala3e097zm7wvvjedzsm2hdc2lsrp",
		Amount:  math.NewInt(44700000),
		Claimed: false,
	},
	{
		Address: "elys17flltxc0vyl7tae3crqe7wvz5ehm4deadqhptv",
		Amount:  math.NewInt(44400000),
		Claimed: false,
	},
	{
		Address: "elys1mmhag8p74eyspn2up33j2zj2gsn45x4fcv6c8a",
		Amount:  math.NewInt(44400000),
		Claimed: false,
	},
	{
		Address: "elys1zf20t43ew6a9fk7g4e4nn5k6w4t3kty87w7pyr",
		Amount:  math.NewInt(44400000),
		Claimed: false,
	},
	{
		Address: "elys148lr9d389eh85puh9e4dvyt40dua62n8uzq4g2",
		Amount:  math.NewInt(44300000),
		Claimed: false,
	},
	{
		Address: "elys189pzd9vxjvtvtwu0pfw3druv9mpq899q4h9u6r",
		Amount:  math.NewInt(44100000),
		Claimed: false,
	},
	{
		Address: "elys19l5r4g7p880efz3tcw3cfen9w86tplffrzlrd0",
		Amount:  math.NewInt(44100000),
		Claimed: false,
	},
	{
		Address: "elys1pmmw2j7q9yh5mm08fmfmtsucx9axdeusxun7h9",
		Amount:  math.NewInt(44100000),
		Claimed: false,
	},
	{
		Address: "elys1qxsel9wjc06jpzpj90dfhvkefssfnwhjzecx0u",
		Amount:  math.NewInt(44100000),
		Claimed: false,
	},
	{
		Address: "elys19m87384xpdylnztrqtuap59gh2dhwee3jjejzj",
		Amount:  math.NewInt(44000000),
		Claimed: false,
	},
	{
		Address: "elys1plynttvhulryjv8cdej5r4ue4zxx82d60vfh4z",
		Amount:  math.NewInt(44000000),
		Claimed: false,
	},
	{
		Address: "elys1yrwh4z7mu6xanywaeejndjlrftjdz6lr5dj5c8",
		Amount:  math.NewInt(44000000),
		Claimed: false,
	},
	{
		Address: "elys1g0cr326krpvfkuv96n7yw5wsyv0y7a2u5c29xe",
		Amount:  math.NewInt(43900000),
		Claimed: false,
	},
	{
		Address: "elys1x2ngs4p3w3jeflqu2unhgtghc3xkjt0pk9wzv0",
		Amount:  math.NewInt(43900000),
		Claimed: false,
	},
	{
		Address: "elys1vydu53mysdnhjtl22dzjz9r36usuvxafl6sqmc",
		Amount:  math.NewInt(43800000),
		Claimed: false,
	},
	{
		Address: "elys10m364n5v7saktreevnwlqz7mz38x8sv2jg9qm0",
		Amount:  math.NewInt(43700000),
		Claimed: false,
	},
	{
		Address: "elys14knz0j70sup7qr5l9588djgaqsvan77dhv86lm",
		Amount:  math.NewInt(43700000),
		Claimed: false,
	},
	{
		Address: "elys190x2tzh5y3l8gtuszdkne7nv3rwkad5xmaq0h4",
		Amount:  math.NewInt(43700000),
		Claimed: false,
	},
	{
		Address: "elys1gpu2ffman85g5cd2zfnwfrp43l5w5kh093q8xl",
		Amount:  math.NewInt(43700000),
		Claimed: false,
	},
	{
		Address: "elys16aszemey8esn4cauv2379qu8a2fz9pdc5frmy6",
		Amount:  math.NewInt(43600000),
		Claimed: false,
	},
	{
		Address: "elys1mntsgfk7gmnywdr79lnqcjzn24eyec9a3jlfvd",
		Amount:  math.NewInt(43600000),
		Claimed: false,
	},
	{
		Address: "elys1w0j3nyqklr6hgkxwz9d7clk9lhkw782y8j8qd6",
		Amount:  math.NewInt(43600000),
		Claimed: false,
	},
	{
		Address: "elys1kp5vq2vqp3x8pygvedspenut075hx89x7nussj",
		Amount:  math.NewInt(43400000),
		Claimed: false,
	},
	{
		Address: "elys1lrw8ghl2ru25pznlyq05jketcgcphtxa6u6j2w",
		Amount:  math.NewInt(43400000),
		Claimed: false,
	},
	{
		Address: "elys1s2xvm5d6kv2mhnwuc4d3xpunqjq2leh9k5xwuz",
		Amount:  math.NewInt(43400000),
		Claimed: false,
	},
	{
		Address: "elys1kfn0g0fpt6z4h4aevqeddh5ul0hk5prrf8x7n6",
		Amount:  math.NewInt(43300000),
		Claimed: false,
	},
	{
		Address: "elys1vy5hzef2k897l4rxkxepjk7tw8hj9fe85y56fh",
		Amount:  math.NewInt(43300000),
		Claimed: false,
	},
	{
		Address: "elys12kek0rvs8rf9nehh5uxeg024zru0j542czt8dq",
		Amount:  math.NewInt(43200000),
		Claimed: false,
	},
	{
		Address: "elys1l57qudtuueg8egm20t0x45g5pwxfkmmtfrj2kk",
		Amount:  math.NewInt(43200000),
		Claimed: false,
	},
	{
		Address: "elys1yk8dsn62an3p8c67a047qv9lvqa7crpawkd3cw",
		Amount:  math.NewInt(43200000),
		Claimed: false,
	},
	{
		Address: "elys1xk0nqswytkdufzwahqnpkyev3w34dedp5shegd",
		Amount:  math.NewInt(43000000),
		Claimed: false,
	},
	{
		Address: "elys1hu00f6fqg0jwf7e9wuqrtv5vh79ttcudwj8akc",
		Amount:  math.NewInt(42800000),
		Claimed: false,
	},
	{
		Address: "elys1ncfq9urxn2rgff2wqsxj645vl7cjkw7fpsgvar",
		Amount:  math.NewInt(42800000),
		Claimed: false,
	},
	{
		Address: "elys18232aj7722efcgmdve4w32clwg99lffvgpdvna",
		Amount:  math.NewInt(42700000),
		Claimed: false,
	},
	{
		Address: "elys1h3kthk2tekp6kkvua7u77m6kyj23fhlpswa8jr",
		Amount:  math.NewInt(42600000),
		Claimed: false,
	},
	{
		Address: "elys1u9en2dyrvg2gxe8se9ktxf8duh6meg2ze287ek",
		Amount:  math.NewInt(42300000),
		Claimed: false,
	},
	{
		Address: "elys1t4zjugxdhwhjss6w42qagc64f7gz7v9t8wetje",
		Amount:  math.NewInt(42200000),
		Claimed: false,
	},
	{
		Address: "elys1za8zv8fd2e9h9qj6tegdx3zqnmmqr3ngwyn8h3",
		Amount:  math.NewInt(42200000),
		Claimed: false,
	},
	{
		Address: "elys17yrl86ehlhuzx2legtfuur7x2wnzvdch6n2jj0",
		Amount:  math.NewInt(42100000),
		Claimed: false,
	},
	{
		Address: "elys1dsfnmlld5hv6vrqhj0dvtes6mlshpju5evytlw",
		Amount:  math.NewInt(42100000),
		Claimed: false,
	},
	{
		Address: "elys1vxrjpcgla4sk6mp3p4jyltp3v08wnplzhlchel",
		Amount:  math.NewInt(42100000),
		Claimed: false,
	},
	{
		Address: "elys180gscdm7ekkwcygwjk808l53ctg2vww07sqae2",
		Amount:  math.NewInt(41900000),
		Claimed: false,
	},
	{
		Address: "elys1gz8uhj09sln2jjlwpe62cnxrw79tzcrhkr02uy",
		Amount:  math.NewInt(41900000),
		Claimed: false,
	},
	{
		Address: "elys1glw20lq2g8vpptxpph6adhp2qpdyzdgzr92upd",
		Amount:  math.NewInt(41800000),
		Claimed: false,
	},
	{
		Address: "elys1p5rqpj64wlhqlchugtk0uv709r6yapazppkkyf",
		Amount:  math.NewInt(41800000),
		Claimed: false,
	},
	{
		Address: "elys1aagtvlny9w2jdldjk2qfp8gu8rtq9uhvtwfkuu",
		Amount:  math.NewInt(41700000),
		Claimed: false,
	},
	{
		Address: "elys1lvd86sr5h8k8gtunj05ct9ajuvmu26sw4av775",
		Amount:  math.NewInt(41700000),
		Claimed: false,
	},
	{
		Address: "elys1v9knk3ymvfvls5kddf0uduvau5lfnzehhy3xlm",
		Amount:  math.NewInt(41700000),
		Claimed: false,
	},
	{
		Address: "elys152azv32wtx9hqys5q5fnfkdy286va4upa0a7dj",
		Amount:  math.NewInt(41600000),
		Claimed: false,
	},
	{
		Address: "elys12yka2w47tx53uhr3az7kr2qxn8u3utq4pdgakf",
		Amount:  math.NewInt(41500000),
		Claimed: false,
	},
	{
		Address: "elys14xwtvm26wg7aenzfkdx5d9q7dg06adp0wzlpn0",
		Amount:  math.NewInt(41500000),
		Claimed: false,
	},
	{
		Address: "elys1swdlr8qsqgzrqe5tpde2jg4cfgzkzwkymhcvrk",
		Amount:  math.NewInt(41500000),
		Claimed: false,
	},
	{
		Address: "elys1w9g8jda6lgnzpsv9f5nr98jvktr0mwgaax83me",
		Amount:  math.NewInt(41500000),
		Claimed: false,
	},
	{
		Address: "elys14zn40x2mm5fmfv6tlcvsapm369g73mz989j3l5",
		Amount:  math.NewInt(41400000),
		Claimed: false,
	},
	{
		Address: "elys1cpny5xfmxwq7ew5u3u0rskcvutt34e2yg5snzy",
		Amount:  math.NewInt(41300000),
		Claimed: false,
	},
	{
		Address: "elys1kpgm2sh48xav7vyvya6s4qc0f9shx9mg8ys3tr",
		Amount:  math.NewInt(41300000),
		Claimed: false,
	},
	{
		Address: "elys1xexzcz2t78s0fmq4qv2ezhfygphlxa4mll4ysg",
		Amount:  math.NewInt(41300000),
		Claimed: false,
	},
	{
		Address: "elys1xvtzkmq5ty5mp75p0rptd0yzlld929jnklrc7t",
		Amount:  math.NewInt(41300000),
		Claimed: false,
	},
	{
		Address: "elys15c2d3dxujcf7my64k0n0hk88g5skwpk3x3l3d0",
		Amount:  math.NewInt(41200000),
		Claimed: false,
	},
	{
		Address: "elys19lw48afwrpasgds4kzhyfn65yhcrfqvft06gj5",
		Amount:  math.NewInt(41200000),
		Claimed: false,
	},
	{
		Address: "elys1pnev36ypmxx4frhc4qj7ex6fe4rz7p74nanv5m",
		Amount:  math.NewInt(41200000),
		Claimed: false,
	},
	{
		Address: "elys1u7eq237d6fg0es765aevmg5874jcxn5g5y3uk4",
		Amount:  math.NewInt(41200000),
		Claimed: false,
	},
	{
		Address: "elys1ucr38czhvm2hcf7gzpcphk7pnxvv9rlx25pw5x",
		Amount:  math.NewInt(41200000),
		Claimed: false,
	},
	{
		Address: "elys1ytmwk4qwphevacssnxaq2c773s0smpmruhvlgp",
		Amount:  math.NewInt(41200000),
		Claimed: false,
	},
	{
		Address: "elys165g9hf24lqdukq6rxndqwmg06ws64u5pxsaykg",
		Amount:  math.NewInt(41100000),
		Claimed: false,
	},
	{
		Address: "elys1fpcuflgyygzwma66zkxy700z80crdrjg7vc7cr",
		Amount:  math.NewInt(41100000),
		Claimed: false,
	},
	{
		Address: "elys18yml5rvvgygudg0lz5xk59vmnsd79kcgn9m5h9",
		Amount:  math.NewInt(41000000),
		Claimed: false,
	},
	{
		Address: "elys1e46svnmq9m03q8284wvf77cggqm55z6f5dpfzf",
		Amount:  math.NewInt(41000000),
		Claimed: false,
	},
	{
		Address: "elys1e7fl3d5u6nrqg8mmkauda7q20l63gnjgg8q7d5",
		Amount:  math.NewInt(41000000),
		Claimed: false,
	},
	{
		Address: "elys1hwjstam766qajmvpuzf2q6rxvpgphypfzp48n2",
		Amount:  math.NewInt(41000000),
		Claimed: false,
	},
	{
		Address: "elys1xmf5dd2p6r68xm800vrksw52qz9r3r4g5tcyz8",
		Amount:  math.NewInt(41000000),
		Claimed: false,
	},
	{
		Address: "elys1048d3468feavzg0j0qzm9xdlh65t52s0znfw8u",
		Amount:  math.NewInt(40900000),
		Claimed: false,
	},
	{
		Address: "elys1973qhhswmsf0y99h88t40v7ecgzqhs6rl4jf50",
		Amount:  math.NewInt(40900000),
		Claimed: false,
	},
	{
		Address: "elys19gczt0nma0c7c0gjfa54dmnrh97mqhmavneqgk",
		Amount:  math.NewInt(40900000),
		Claimed: false,
	},
	{
		Address: "elys1j6hwm5knyhxs8fdvge36tne57xeqsk0e9w9sgq",
		Amount:  math.NewInt(40900000),
		Claimed: false,
	},
	{
		Address: "elys1m9jpccyg6npez0nwueszd6rxrrk0r2gwl8kaq5",
		Amount:  math.NewInt(40900000),
		Claimed: false,
	},
	{
		Address: "elys1w7ln88mjzsadhj237evy7rejaufvfavlnvwjn8",
		Amount:  math.NewInt(40900000),
		Claimed: false,
	},
	{
		Address: "elys135nzejn46u752u7uqkdd3r35jpa23q9mg7aczz",
		Amount:  math.NewInt(40800000),
		Claimed: false,
	},
	{
		Address: "elys1hewyn5hpgzcngy03ztgkyp4mv3w5nar6xkcq55",
		Amount:  math.NewInt(40800000),
		Claimed: false,
	},
	{
		Address: "elys1zq404dzfgcf077aylref0ld5jst456se3w936f",
		Amount:  math.NewInt(40800000),
		Claimed: false,
	},
	{
		Address: "elys1w0jv82w3j7a985px454fuuuhdlv90uxq9ptg4t",
		Amount:  math.NewInt(40700000),
		Claimed: false,
	},
	{
		Address: "elys1wrkp5fdhp60236nekax2tqhux2d3c9l0rvxddr",
		Amount:  math.NewInt(40700000),
		Claimed: false,
	},
	{
		Address: "elys1j36yjklwcaunaaszv9e52p7svmje457h7z8vrp",
		Amount:  math.NewInt(40500000),
		Claimed: false,
	},
	{
		Address: "elys1p4f93qkjghdkp54rtvdg5vpzmsgcep5pvgpuf8",
		Amount:  math.NewInt(40500000),
		Claimed: false,
	},
	{
		Address: "elys1pylujqk4e4leq5vuly9t6tu78c9md54slkqasy",
		Amount:  math.NewInt(40500000),
		Claimed: false,
	},
	{
		Address: "elys1jkpsqlxguk535kp9g2p0kfdus63cql2plwj796",
		Amount:  math.NewInt(40400000),
		Claimed: false,
	},
	{
		Address: "elys15t2gn2ygna49uqal4lu2hfwc8g3c0erp25ttf7",
		Amount:  math.NewInt(40200000),
		Claimed: false,
	},
	{
		Address: "elys15xjk5j65dad7zmgwh5hx5re087asxyq4hdye8w",
		Amount:  math.NewInt(40200000),
		Claimed: false,
	},
	{
		Address: "elys1r5trk7akrgukgcvvf852kwzgzydmla3rde493a",
		Amount:  math.NewInt(40200000),
		Claimed: false,
	},
	{
		Address: "elys1ktvnuysk4lz0westx7wg9q6rggz8e59zfmaycl",
		Amount:  math.NewInt(40100000),
		Claimed: false,
	},
	{
		Address: "elys1um5ykppchhzj7cu0ag2h4d3ux2mhruc9h6z9qw",
		Amount:  math.NewInt(40100000),
		Claimed: false,
	},
	{
		Address: "elys1z5437s7v792lepppt729taz3369kjelmtsr53r",
		Amount:  math.NewInt(40100000),
		Claimed: false,
	},
	{
		Address: "elys142j47nzzxzesmngfp4el4sp9682907urmgjnkh",
		Amount:  math.NewInt(40000000),
		Claimed: false,
	},
	{
		Address: "elys13u5tllqrq2u6vslkdqxrywe7mrjnfc7r9vlysj",
		Amount:  math.NewInt(39900000),
		Claimed: false,
	},
	{
		Address: "elys19l43sc0p8455whqj4h0td5m08ru99ze0hxsurl",
		Amount:  math.NewInt(39900000),
		Claimed: false,
	},
	{
		Address: "elys1e2rnuchylkjku6ztxzstszhca578cvmymp2mf9",
		Amount:  math.NewInt(39900000),
		Claimed: false,
	},
	{
		Address: "elys1gkmp5gvzwt6wfru4uajzqmgj0gv0lzq3alsv0m",
		Amount:  math.NewInt(39900000),
		Claimed: false,
	},
	{
		Address: "elys1uwsypxpdjpdfcmfh5nlu4sj6f9y53x65d5aulu",
		Amount:  math.NewInt(39900000),
		Claimed: false,
	},
	{
		Address: "elys1et73jw4lqnjp488vd9m423762yjh6e3lkljev9",
		Amount:  math.NewInt(39800000),
		Claimed: false,
	},
	{
		Address: "elys1jcqtsze0gda7lrr2a2qm86twrhdj29caq3tcgh",
		Amount:  math.NewInt(39800000),
		Claimed: false,
	},
	{
		Address: "elys1rj5yn8mg36emyp79pr25dkqymwpqc76f5dw5r4",
		Amount:  math.NewInt(39800000),
		Claimed: false,
	},
	{
		Address: "elys1c80x62uj8y4qq2zclfm5533qjlpkqqljf9tas4",
		Amount:  math.NewInt(39700000),
		Claimed: false,
	},
	{
		Address: "elys1r0y3azermt32tjemw5ns87a8mvnn3tqag3z8uu",
		Amount:  math.NewInt(39700000),
		Claimed: false,
	},
	{
		Address: "elys19teekcma2wugx08yrp2wwv98vgtfx6rmytaax3",
		Amount:  math.NewInt(39600000),
		Claimed: false,
	},
	{
		Address: "elys1dneehk7egpzvwfp4cuhcwgxa9zsn6dwczk9fxm",
		Amount:  math.NewInt(39600000),
		Claimed: false,
	},
	{
		Address: "elys1ha9qswmzk6s9z90x0y9vclmkzmx6s8h7j8d3zs",
		Amount:  math.NewInt(39600000),
		Claimed: false,
	},
	{
		Address: "elys1pwp8p47hxjfgkt0fk2m4x2avgx7shyjzhrh027",
		Amount:  math.NewInt(39600000),
		Claimed: false,
	},
	{
		Address: "elys1yjdhjjkwmwqgskwvwtst4n0tnynkzfl8jyp6cp",
		Amount:  math.NewInt(39600000),
		Claimed: false,
	},
	{
		Address: "elys18yfg6kvrvem3e4zle506due53t9du85du7cf3u",
		Amount:  math.NewInt(39500000),
		Claimed: false,
	},
	{
		Address: "elys1j7yptwyynr89hhl9c3yzf4ks6ehhcekzuhk5ya",
		Amount:  math.NewInt(39500000),
		Claimed: false,
	},
	{
		Address: "elys1lzhwkl3388mllfrhuzptqlsavsd39l0332603g",
		Amount:  math.NewInt(39400000),
		Claimed: false,
	},
	{
		Address: "elys16eexps7gmn28x7y734uyx5578d5u395nnmsxnx",
		Amount:  math.NewInt(39300000),
		Claimed: false,
	},
	{
		Address: "elys17r9uvhzmqrg7t490kfmq8eqa7dw8lh7upwe45d",
		Amount:  math.NewInt(39300000),
		Claimed: false,
	},
	{
		Address: "elys1dtyj5csdnpn0rsnrm9dqhyj2aapvu7l5clzkry",
		Amount:  math.NewInt(39200000),
		Claimed: false,
	},
	{
		Address: "elys1nhvqlrvglxq38uj9d4fwqpdggj7u4486p0ry4l",
		Amount:  math.NewInt(39200000),
		Claimed: false,
	},
	{
		Address: "elys15suc5lfnq05vjauumalfgkayxa6s9lsva2770k",
		Amount:  math.NewInt(39100000),
		Claimed: false,
	},
	{
		Address: "elys17xnxpzw243mj2h2k67xty3tgykcu2hwfg0q798",
		Amount:  math.NewInt(39100000),
		Claimed: false,
	},
	{
		Address: "elys1t9gwdn60wx7cnwhgzr659rspd508gexaypkjmy",
		Amount:  math.NewInt(39100000),
		Claimed: false,
	},
	{
		Address: "elys1xgx7xwd8wxrumd0mz6zhdzk6j4s5ycuq86xwyc",
		Amount:  math.NewInt(39100000),
		Claimed: false,
	},
	{
		Address: "elys1xh8ydf9rqxtyg7txh56wrcysadk3htflgna0vs",
		Amount:  math.NewInt(39100000),
		Claimed: false,
	},
	{
		Address: "elys1cscup7vc2tjcjlx2cq30d2uknuegxy0ptvktmx",
		Amount:  math.NewInt(39000000),
		Claimed: false,
	},
	{
		Address: "elys1mpw65f5l2memamjagh7a35p747dhe8jm7usr2f",
		Amount:  math.NewInt(39000000),
		Claimed: false,
	},
	{
		Address: "elys1u28n46lwz2j0sjsagk0wyum5aqfwck02gykwxq",
		Amount:  math.NewInt(39000000),
		Claimed: false,
	},
	{
		Address: "elys1yg36eryq2zgq23wnhw3fjs4ax45slw7ay76tgf",
		Amount:  math.NewInt(39000000),
		Claimed: false,
	},
	{
		Address: "elys1clcd6n3uj58efwapv7sg5uwtdzexasxuvj9aqz",
		Amount:  math.NewInt(38800000),
		Claimed: false,
	},
	{
		Address: "elys1z8hj80cj2lap7e4kadsg4w98ytgpxs6x9azkxl",
		Amount:  math.NewInt(38800000),
		Claimed: false,
	},
	{
		Address: "elys1zcps8slrrmajq3la7u8eaem4sv0vtx4yq726qe",
		Amount:  math.NewInt(38800000),
		Claimed: false,
	},
	{
		Address: "elys1442jxc2sdh0rpw20lt0w2gep50kypjy3vq3u7e",
		Amount:  math.NewInt(38700000),
		Claimed: false,
	},
	{
		Address: "elys1jne5k8jzvdp8rkdjf5h9x473528du6arv5jk55",
		Amount:  math.NewInt(38700000),
		Claimed: false,
	},
	{
		Address: "elys1w5t033pvqgs44g8lnalh8j2az96t76kfckgqw4",
		Amount:  math.NewInt(38700000),
		Claimed: false,
	},
	{
		Address: "elys15uwq74308cr6q5u4mxq5spge2ucey8zylk3tx9",
		Amount:  math.NewInt(38600000),
		Claimed: false,
	},
	{
		Address: "elys1cklntpzz5syyg2xadxvty49zhcuwayxxrcqmln",
		Amount:  math.NewInt(38600000),
		Claimed: false,
	},
	{
		Address: "elys1frag0anh77kvr9mlpp9g3dta40cv3vgy60mj2f",
		Amount:  math.NewInt(38600000),
		Claimed: false,
	},
	{
		Address: "elys1ptqdzjscf2986pc3ndz5cr6xj5p5ca2avm6p0k",
		Amount:  math.NewInt(38600000),
		Claimed: false,
	},
	{
		Address: "elys14uzu9dcdqaw0fnnxtjk0uy92c3rv3ae6y6e4ql",
		Amount:  math.NewInt(38500000),
		Claimed: false,
	},
	{
		Address: "elys1lrptcw7knr46yzkns0t09t7ny0yurwk5kqu520",
		Amount:  math.NewInt(38400000),
		Claimed: false,
	},
	{
		Address: "elys1sl0wrl5ntjd7qfwvltz42j875vu3yr64773a5m",
		Amount:  math.NewInt(38400000),
		Claimed: false,
	},
	{
		Address: "elys1x2qrp5p6t27leyw7hcql2fzs0q56r6uv0d8nhk",
		Amount:  math.NewInt(38400000),
		Claimed: false,
	},
	{
		Address: "elys1yya5f93yjqkj5d2ncayxz5lzzu3t524rnpug7v",
		Amount:  math.NewInt(38400000),
		Claimed: false,
	},
	{
		Address: "elys1cv7n60quu7d97gt9krz7h9ctvvsep6q6phw0rd",
		Amount:  math.NewInt(38300000),
		Claimed: false,
	},
	{
		Address: "elys10e7qugnq9jrxqq2axh93ecvsmmeydfaxutqfcd",
		Amount:  math.NewInt(38200000),
		Claimed: false,
	},
	{
		Address: "elys1h3zkdsw2gth67ft4yl70cqvlha60jnxq64sx9u",
		Amount:  math.NewInt(38200000),
		Claimed: false,
	},
	{
		Address: "elys1h9dudev2tux7nayy4xqeyzpwrg5e7f9yvqv8ve",
		Amount:  math.NewInt(38200000),
		Claimed: false,
	},
	{
		Address: "elys1ltxfradxpfw42v0t6wls8ypkcsvc9kypdr5hf2",
		Amount:  math.NewInt(38200000),
		Claimed: false,
	},
	{
		Address: "elys1zqf8mze8ljxcqp38fqs9ceynck277e6m072zcn",
		Amount:  math.NewInt(38200000),
		Claimed: false,
	},
	{
		Address: "elys15299fs3sdypvwm5d9src0hcv7cvn8fuhc9rh4m",
		Amount:  math.NewInt(38100000),
		Claimed: false,
	},
	{
		Address: "elys1ckwlqjx8nr5j9qy5as0hje03l2hw9cpsrcgmpw",
		Amount:  math.NewInt(38100000),
		Claimed: false,
	},
	{
		Address: "elys1ud0t3z69a04csu9gcceg4vrvt0maf9v9vdvsjy",
		Amount:  math.NewInt(38100000),
		Claimed: false,
	},
	{
		Address: "elys17h6wp5ewmnet6xxsa48yr572zpya0kqz9wv8df",
		Amount:  math.NewInt(38000000),
		Claimed: false,
	},
	{
		Address: "elys1qem7lgkey4j4w9kvtvvn78g5cnaqehvysx04jh",
		Amount:  math.NewInt(38000000),
		Claimed: false,
	},
	{
		Address: "elys1wdvj5sf55cekusumjfg4sfdz2vzwjhnu6lz65t",
		Amount:  math.NewInt(38000000),
		Claimed: false,
	},
	{
		Address: "elys1k4paygdg7tyw3ye65qua00pw6vlqmvsk0qyym5",
		Amount:  math.NewInt(37800000),
		Claimed: false,
	},
	{
		Address: "elys1m4k96fm24xksjkytcx97vwwpf27qj26xlhpfj5",
		Amount:  math.NewInt(37800000),
		Claimed: false,
	},
	{
		Address: "elys1qudhml67zaqlepv4js97sj4a83483t9xj6skag",
		Amount:  math.NewInt(37800000),
		Claimed: false,
	},
	{
		Address: "elys1w28my53zxnh8hk043vrkhpct7ew78m8gvuyakf",
		Amount:  math.NewInt(37800000),
		Claimed: false,
	},
	{
		Address: "elys1yrjq3ya5njp6n0sn4kt3urgwhlduxe5jtrhe2l",
		Amount:  math.NewInt(37800000),
		Claimed: false,
	},
	{
		Address: "elys1s9mzdutyj295cwjcdrk6wnd26wfm2dnw04qweh",
		Amount:  math.NewInt(37700000),
		Claimed: false,
	},
	{
		Address: "elys1ys6hxfh6w956uj99vcgrf0s9h2tw6a4rtf0tak",
		Amount:  math.NewInt(37700000),
		Claimed: false,
	},
	{
		Address: "elys19k3l0z5txzv2w3kr28lhs49dfd65sv63zw8gy0",
		Amount:  math.NewInt(37600000),
		Claimed: false,
	},
	{
		Address: "elys1tp76707mftgqucpez8p6ff0ypkvu85tghrjke7",
		Amount:  math.NewInt(37600000),
		Claimed: false,
	},
	{
		Address: "elys186l3zrq6qh20hunyycv0cuu242g2yt3atu2gvt",
		Amount:  math.NewInt(37500000),
		Claimed: false,
	},
	{
		Address: "elys1fcfpmgmqkyyrtag7z3qdvezkfgpmumq8estenu",
		Amount:  math.NewInt(37500000),
		Claimed: false,
	},
	{
		Address: "elys1myw2m6ty2es7qv8aj7m74szv4ju5yndnam0u43",
		Amount:  math.NewInt(37500000),
		Claimed: false,
	},
	{
		Address: "elys148rfszw49j4yng9hm30t9fx3ptf30wmg4hhwnn",
		Amount:  math.NewInt(37400000),
		Claimed: false,
	},
	{
		Address: "elys1494tc0fn8euxx0z2yyehjj9myn2k5xsxjqzy8a",
		Amount:  math.NewInt(37400000),
		Claimed: false,
	},
	{
		Address: "elys186mgcxr4mgapa5pkvgdff3qukql030ceg050xh",
		Amount:  math.NewInt(37400000),
		Claimed: false,
	},
	{
		Address: "elys18uht24fddp9px6plg7llw0jv9dqz3m96yu0dl9",
		Amount:  math.NewInt(37400000),
		Claimed: false,
	},
	{
		Address: "elys1wujxr8a3tg232jnw4an42szs62adnvt5u6w2eu",
		Amount:  math.NewInt(37400000),
		Claimed: false,
	},
	{
		Address: "elys195na5da0hn20wl4qz5gh9f5xttsdw08yhhl82z",
		Amount:  math.NewInt(37300000),
		Claimed: false,
	},
	{
		Address: "elys1m9ntttsjxhqwydlhsg550uv4twp20gu4cjzl9l",
		Amount:  math.NewInt(37300000),
		Claimed: false,
	},
	{
		Address: "elys1nf8zr9uglvzyv575kwn3pgkexw9rzu7dc8mse2",
		Amount:  math.NewInt(37300000),
		Claimed: false,
	},
	{
		Address: "elys168dk0qv4a62t0rmr7ykw2atd77aeca5dt8z45q",
		Amount:  math.NewInt(37200000),
		Claimed: false,
	},
	{
		Address: "elys1gjlcn867q897yfpmdvlhpx9umn3y8l7ah08srh",
		Amount:  math.NewInt(37200000),
		Claimed: false,
	},
	{
		Address: "elys1jtt8wa3tr8v5j45madamzxu5ysw0nqu5m2zxmh",
		Amount:  math.NewInt(37200000),
		Claimed: false,
	},
	{
		Address: "elys1kh9htzc64hsydna8pprmj2vpqd6524xh2gazy8",
		Amount:  math.NewInt(37200000),
		Claimed: false,
	},
	{
		Address: "elys1auw85vqffumwuljhz3924p3n48960ztghjl6fc",
		Amount:  math.NewInt(37100000),
		Claimed: false,
	},
	{
		Address: "elys1d35pgdnt80s3ydclrqw38emzkxsh7vj50g9yqv",
		Amount:  math.NewInt(37100000),
		Claimed: false,
	},
	{
		Address: "elys1fdm6prm6nn9snm7vgj943fz0fgj70xdmcxk427",
		Amount:  math.NewInt(37100000),
		Claimed: false,
	},
	{
		Address: "elys1kxh247kpunurpuxykydqhgz96jelrvlqamy7ql",
		Amount:  math.NewInt(37100000),
		Claimed: false,
	},
	{
		Address: "elys1nhs4jkz6k2fwnh4ft7ea95y3hfyu8fvp4j7cet",
		Amount:  math.NewInt(37100000),
		Claimed: false,
	},
	{
		Address: "elys14emp06x7lkd7zyzuylyej9ag2ahsjmp0u63swc",
		Amount:  math.NewInt(37000000),
		Claimed: false,
	},
	{
		Address: "elys1f7pp2w6lv82un7zc00susjatepyus4cffej83d",
		Amount:  math.NewInt(37000000),
		Claimed: false,
	},
	{
		Address: "elys1szy5qhnmvcv4a9psfz75ezh8updf6rrrgahq42",
		Amount:  math.NewInt(37000000),
		Claimed: false,
	},
	{
		Address: "elys1z8zvrpmmv5l3punazazkesdgqeev098sj0397m",
		Amount:  math.NewInt(37000000),
		Claimed: false,
	},
	{
		Address: "elys126uhl9qvf0qdualqdpz6lqaptn8x4kuus5es7e",
		Amount:  math.NewInt(36800000),
		Claimed: false,
	},
	{
		Address: "elys13u8haeq8v2cuyvvy8ya3t5ynzpve6rm89508ru",
		Amount:  math.NewInt(36800000),
		Claimed: false,
	},
	{
		Address: "elys1pzt9lmpdwuaqsdrter5k03jf4hss8gy9nds6j8",
		Amount:  math.NewInt(36700000),
		Claimed: false,
	},
	{
		Address: "elys1v6sw9s9j669rfqlsyg77cp46a3g08zyjpqe7al",
		Amount:  math.NewInt(36700000),
		Claimed: false,
	},
	{
		Address: "elys1wg6uw7t05dj0akjqvz54q4ydqa62tu40rzwl38",
		Amount:  math.NewInt(36700000),
		Claimed: false,
	},
	{
		Address: "elys1kqh9v54k9qvrz5ggplndzlgynflk2n33pl4tut",
		Amount:  math.NewInt(36600000),
		Claimed: false,
	},
	{
		Address: "elys1s4u6wh0jw4a5nesm593atlvljgvddmczg6lxz9",
		Amount:  math.NewInt(36600000),
		Claimed: false,
	},
	{
		Address: "elys1uakeurewqzxmhaktmhunfk8js2ms5p4hvgnspl",
		Amount:  math.NewInt(36600000),
		Claimed: false,
	},
	{
		Address: "elys1vqv5sug9s3sg54uq3pq5u8nl2uwgeuwrfvm2ru",
		Amount:  math.NewInt(36600000),
		Claimed: false,
	},
	{
		Address: "elys1yzte5q5qj7l9uldu6p8k32536yg9tp47sqm8c6",
		Amount:  math.NewInt(36500000),
		Claimed: false,
	},
	{
		Address: "elys14p32amp5cdaq2t5cdacnsd92gcdykpdvjp7n3p",
		Amount:  math.NewInt(36400000),
		Claimed: false,
	},
	{
		Address: "elys14qupf2xsltxgp3w3uyfynl0qqdcvva3057q705",
		Amount:  math.NewInt(36400000),
		Claimed: false,
	},
	{
		Address: "elys1ag7p5y4k736l4dctjrn2c6ym72rjqszce3w09c",
		Amount:  math.NewInt(36400000),
		Claimed: false,
	},
	{
		Address: "elys1eq8q2zzjg77rsdfmdhyt2pu6lcqxvucg0s5feh",
		Amount:  math.NewInt(36400000),
		Claimed: false,
	},
	{
		Address: "elys1xrcwyvxcfm0jftp4zlmhzcj33pwp3cu443cvtr",
		Amount:  math.NewInt(36400000),
		Claimed: false,
	},
	{
		Address: "elys13jc3cvnpzgetp67ze8kc068x6ng7atezp4x7ug",
		Amount:  math.NewInt(36200000),
		Claimed: false,
	},
	{
		Address: "elys1uv0k22f64pysjy4pnelfjz2nnf8k4nmwtmsl85",
		Amount:  math.NewInt(36200000),
		Claimed: false,
	},
	{
		Address: "elys19x3s23x9hvp98rdgrl6mnec8h5hy7zmnn6c324",
		Amount:  math.NewInt(36100000),
		Claimed: false,
	},
	{
		Address: "elys1fy8s6yyn58uwhvd3psxyku72mycgtl6rnh63cg",
		Amount:  math.NewInt(36100000),
		Claimed: false,
	},
	{
		Address: "elys10kgxlyyl3cndaxa62x64alczzpzkqqg2vadf23",
		Amount:  math.NewInt(36000000),
		Claimed: false,
	},
	{
		Address: "elys17d87xpt4f7rp75q2y36y659k73m9n9fjtuz3f6",
		Amount:  math.NewInt(36000000),
		Claimed: false,
	},
	{
		Address: "elys1memcf6sqd00h6yzf0e8e8e707tjmszpdq420dk",
		Amount:  math.NewInt(36000000),
		Claimed: false,
	},
	{
		Address: "elys157tdzyvj0vce2vrg8q5sep7rvamnf2tqfxsl9l",
		Amount:  math.NewInt(35900000),
		Claimed: false,
	},
	{
		Address: "elys17hs4un85e5nqgf3l2ufe967ppalvzzphnkqpv7",
		Amount:  math.NewInt(35900000),
		Claimed: false,
	},
	{
		Address: "elys1ceyaz9djhmxh3m3zyhycktxqjd6ezhpnvnd0ac",
		Amount:  math.NewInt(35900000),
		Claimed: false,
	},
	{
		Address: "elys1dswdrk3qfqupwpun29df5vrt62tku2jyah70zk",
		Amount:  math.NewInt(35900000),
		Claimed: false,
	},
	{
		Address: "elys1r54z6xduqcd7de3w8ecn0ylw3x8wa3eaz00yee",
		Amount:  math.NewInt(35900000),
		Claimed: false,
	},
	{
		Address: "elys1xzccphkekmnxekr7ddfr4f2mydgx6cadnkpsmh",
		Amount:  math.NewInt(35900000),
		Claimed: false,
	},
	{
		Address: "elys185y02npy45f9nq9vhzcz4y5822svjfjyueum9j",
		Amount:  math.NewInt(35800000),
		Claimed: false,
	},
	{
		Address: "elys1atvc7s45f7zwm9rh3dqrtd4ued2dknhym8z0wx",
		Amount:  math.NewInt(35800000),
		Claimed: false,
	},
	{
		Address: "elys1xns4z59f6suxmdtc3cf8avejmzr3ykr63gq3eq",
		Amount:  math.NewInt(35800000),
		Claimed: false,
	},
	{
		Address: "elys1fqfkwftgxehk0vzfmn24dmnpt2t8ekedxru3a7",
		Amount:  math.NewInt(35700000),
		Claimed: false,
	},
	{
		Address: "elys1smak362ve6qk3d8grw4um4gcz3frflt0cers8l",
		Amount:  math.NewInt(35700000),
		Claimed: false,
	},
	{
		Address: "elys1nzedewmlnlpwsmu2kh3g7mvhc02wp8k29k5y7q",
		Amount:  math.NewInt(35600000),
		Claimed: false,
	},
	{
		Address: "elys1znaj4lqyds2jsg9tpmwt59r0mxn85at48ggqzh",
		Amount:  math.NewInt(35600000),
		Claimed: false,
	},
	{
		Address: "elys1trkq34l29cespl2llx6j0nxjgng7dau9v96y8q",
		Amount:  math.NewInt(35500000),
		Claimed: false,
	},
	{
		Address: "elys1uddc238v9whrntq7445rx3d8eah53yd7nruz08",
		Amount:  math.NewInt(35400000),
		Claimed: false,
	},
	{
		Address: "elys1uvpx8jcrkn26tmtxeh57gmsfpe406fpnz4vn8v",
		Amount:  math.NewInt(35400000),
		Claimed: false,
	},
	{
		Address: "elys166csq0r6nhvk5eutyhcp4zqpldudwd6t5sfv6x",
		Amount:  math.NewInt(35300000),
		Claimed: false,
	},
	{
		Address: "elys1cka238xe6fpczw8rn645rlpdepwdj80g46v4ym",
		Amount:  math.NewInt(35300000),
		Claimed: false,
	},
	{
		Address: "elys1w22llzhz5d7mqh7paycyak6qt0xdd26qkw9yuq",
		Amount:  math.NewInt(35300000),
		Claimed: false,
	},
	{
		Address: "elys1wdun0av37e76ax7t4egpxndmx84yurkyj6cx0g",
		Amount:  math.NewInt(35300000),
		Claimed: false,
	},
	{
		Address: "elys1p9ytkzydys9779l0swvlyhnaehvnf9ngtrtlv6",
		Amount:  math.NewInt(35200000),
		Claimed: false,
	},
	{
		Address: "elys12fx3ajf6xcwy6j4uzcsw4klcqdfjldedhtjy3x",
		Amount:  math.NewInt(35100000),
		Claimed: false,
	},
	{
		Address: "elys1909mrytp2p2mqkd9xjasyfaurhx5df4t4cfufl",
		Amount:  math.NewInt(35100000),
		Claimed: false,
	},
	{
		Address: "elys14yq6js8kee2uqjkmmusf4wt00r2mj3ju4ps2we",
		Amount:  math.NewInt(34900000),
		Claimed: false,
	},
	{
		Address: "elys1k2ahl9y5q6qx7mwf80plpmu8en7t3et8qlu87t",
		Amount:  math.NewInt(34900000),
		Claimed: false,
	},
	{
		Address: "elys1rjszx2t3m7a6spel2atg67ph6zlrwk2fe6s709",
		Amount:  math.NewInt(34900000),
		Claimed: false,
	},
	{
		Address: "elys126xfrtkj2q5azmud6r8m9h46sxw4vy5xqdmm9n",
		Amount:  math.NewInt(34700000),
		Claimed: false,
	},
	{
		Address: "elys1mjttyf5qwcuepumeq2tjfhxm9t54v3sq2zdcw7",
		Amount:  math.NewInt(34700000),
		Claimed: false,
	},
	{
		Address: "elys1rhf6qxezm8xd3ma2yernm5u2777x2h2kup4ev0",
		Amount:  math.NewInt(34700000),
		Claimed: false,
	},
	{
		Address: "elys1qksdaex4ws9rdakalcm2yl59uxfcxm5dkrxlqd",
		Amount:  math.NewInt(34600000),
		Claimed: false,
	},
	{
		Address: "elys1h3pd599uevns4gfl84ha3yuc4zll5xlzxj8ka7",
		Amount:  math.NewInt(34500000),
		Claimed: false,
	},
	{
		Address: "elys1rhltdhqldnex5fuuk53y9syej885l04ecnmxa5",
		Amount:  math.NewInt(34500000),
		Claimed: false,
	},
	{
		Address: "elys10w0qn69xk9jcy029unu6eakgc6ldyzd25nu560",
		Amount:  math.NewInt(34400000),
		Claimed: false,
	},
	{
		Address: "elys17wlgtuklu6y3a8f6r0xtnv40smq3l7fd2dkjqr",
		Amount:  math.NewInt(34400000),
		Claimed: false,
	},
	{
		Address: "elys1hpkyw50ll59rp7atem8tvsjfaj8f2kpjpxew8h",
		Amount:  math.NewInt(34400000),
		Claimed: false,
	},
	{
		Address: "elys1jks33wuwddqr2ncu2nw6l9zn0y8rmq98c6kgj7",
		Amount:  math.NewInt(34400000),
		Claimed: false,
	},
	{
		Address: "elys1qyj8e0rhgfwxq4tupupdwlazzgkwessp5wx6ru",
		Amount:  math.NewInt(34400000),
		Claimed: false,
	},
	{
		Address: "elys1s6v4f00zd6nf6hjzwy7y3exfwjdc9q5ekr9uh4",
		Amount:  math.NewInt(34400000),
		Claimed: false,
	},
	{
		Address: "elys1alnfxadhuu4x6ftcmwzeks86de4rk56ss76pp0",
		Amount:  math.NewInt(34300000),
		Claimed: false,
	},
	{
		Address: "elys1fdm9hmeganswpasfx7zne2duunah5sxm6ete47",
		Amount:  math.NewInt(34300000),
		Claimed: false,
	},
	{
		Address: "elys1fea9dkekqd6ngjn4nulkwgahtjky3xjy7vz3wx",
		Amount:  math.NewInt(34300000),
		Claimed: false,
	},
	{
		Address: "elys1qvlq79wdt8j7cvexd9jvk30smkn6sqfhh3whqn",
		Amount:  math.NewInt(34300000),
		Claimed: false,
	},
	{
		Address: "elys1rwdut94wp03hew5d7ulwg2m8gg5ww6kkpjxll0",
		Amount:  math.NewInt(34300000),
		Claimed: false,
	},
	{
		Address: "elys1rxjw8rxhghc599cupzp6h6h6sfssxdu6vkph3k",
		Amount:  math.NewInt(34300000),
		Claimed: false,
	},
	{
		Address: "elys1tu8u9gekvrg4c7uh8flh0jn5symmmwvwjke48l",
		Amount:  math.NewInt(34200000),
		Claimed: false,
	},
	{
		Address: "elys1s85dql34tlqa4h92wa52jvv4lxwvp6xlmndwn4",
		Amount:  math.NewInt(34100000),
		Claimed: false,
	},
	{
		Address: "elys12d5s7d69as2y8dl3mgx0ldw2vaj2r879jgcxyn",
		Amount:  math.NewInt(34000000),
		Claimed: false,
	},
	{
		Address: "elys1fth2nyd6g9q03n90ddrljfcxa8vja0wc8negut",
		Amount:  math.NewInt(34000000),
		Claimed: false,
	},
	{
		Address: "elys15d9wr0r5aqecckqlc5v87n88yrvrduwea3h9g7",
		Amount:  math.NewInt(33900000),
		Claimed: false,
	},
	{
		Address: "elys19fm4mqhn4zyejxlythwn7kaqz359d5yf20enp4",
		Amount:  math.NewInt(33900000),
		Claimed: false,
	},
	{
		Address: "elys1m3ygpwf5e8cn49ue9du223y2qlzxtzswuw8qhc",
		Amount:  math.NewInt(33900000),
		Claimed: false,
	},
	{
		Address: "elys1x9rjvs60y8uzqtzgyms7e0c3qfy6ulvkuacpdc",
		Amount:  math.NewInt(33800000),
		Claimed: false,
	},
	{
		Address: "elys1azm58aqs4enp09ndema9ufu8jph77rqvph0xae",
		Amount:  math.NewInt(33600000),
		Claimed: false,
	},
	{
		Address: "elys175rl5xx6k83p7gwnmf9zww0frt776h332vtal8",
		Amount:  math.NewInt(33500000),
		Claimed: false,
	},
	{
		Address: "elys1aenrp6mh2mrv2e6caddvjvuq5ere3htcc67dt3",
		Amount:  math.NewInt(33500000),
		Claimed: false,
	},
	{
		Address: "elys1d27st6a08jau63w0a2e4axfvzlv5ssqlaqys4r",
		Amount:  math.NewInt(33500000),
		Claimed: false,
	},
	{
		Address: "elys1f75hycxudmrtvwm6mn53wqvd2y4k6f7lqrhw08",
		Amount:  math.NewInt(33500000),
		Claimed: false,
	},
	{
		Address: "elys1ucj4xmmq3ykgkw2halpf24h5n2awsu2thr5qnu",
		Amount:  math.NewInt(33500000),
		Claimed: false,
	},
	{
		Address: "elys108aknlpkxy5sn7y0k3szd3cln8hxhfaxgnh52v",
		Amount:  math.NewInt(33400000),
		Claimed: false,
	},
	{
		Address: "elys10f9ffs6enltrfc5kg6np28hwyvtd59smq8wssf",
		Amount:  math.NewInt(33400000),
		Claimed: false,
	},
	{
		Address: "elys18nxfv4qcvn9qz55v268p59aehfkzj9ycsuu2ep",
		Amount:  math.NewInt(33400000),
		Claimed: false,
	},
	{
		Address: "elys1c79q60sq080cljg399synpxgn038hwrkm5hw9a",
		Amount:  math.NewInt(33400000),
		Claimed: false,
	},
	{
		Address: "elys1v9rtrkeufr7yxy0p24wwj4u5rwyalud3326wqa",
		Amount:  math.NewInt(33400000),
		Claimed: false,
	},
	{
		Address: "elys15kaueuu45w47nd3fh8ul3773w0xx7adcqdu6p4",
		Amount:  math.NewInt(33299999),
		Claimed: false,
	},
	{
		Address: "elys1cwl957d87xj4tud0wj9ry9vygympaaapta39sa",
		Amount:  math.NewInt(33299999),
		Claimed: false,
	},
	{
		Address: "elys1qn459sd3wpfz4epjndgermgyvcv9xk0vhsw54y",
		Amount:  math.NewInt(33200000),
		Claimed: false,
	},
	{
		Address: "elys1hrhsrwxyt84q3k3fffc6d8ert2yxhglurkmd79",
		Amount:  math.NewInt(33100000),
		Claimed: false,
	},
	{
		Address: "elys1ku47rda802atptcgunlhv2pjnyktyzznntgwec",
		Amount:  math.NewInt(33100000),
		Claimed: false,
	},
	{
		Address: "elys1taac7pyptemjfuhkjj4ph5qt8hjxe8n5hgq49w",
		Amount:  math.NewInt(33100000),
		Claimed: false,
	},
	{
		Address: "elys1xjrpscg790jlv45edw2ww04fkced8ff2d9a7gh",
		Amount:  math.NewInt(33100000),
		Claimed: false,
	},
	{
		Address: "elys1rg5z6zf9lym94e5qeuk6wdh4gdy54kfj8r294y",
		Amount:  math.NewInt(33000000),
		Claimed: false,
	},
	{
		Address: "elys1s3htt50lud8avwsgmqsvqgjelh6raawuxleeh4",
		Amount:  math.NewInt(33000000),
		Claimed: false,
	},
	{
		Address: "elys12ufxmc72lgkrj4s45mjt87hpnv7d5j5yqdmt97",
		Amount:  math.NewInt(32900000),
		Claimed: false,
	},
	{
		Address: "elys1prc9tnpk7ze3xvcxpz4kud57vgnyjka2k0teft",
		Amount:  math.NewInt(32900000),
		Claimed: false,
	},
	{
		Address: "elys120ccruea2enrgmx2krdd5u8r2y0lnvh30hhs9x",
		Amount:  math.NewInt(32799999),
		Claimed: false,
	},
	{
		Address: "elys1f22jcx3u470vd6h7tl2wvmfz65vdguwtgptz0w",
		Amount:  math.NewInt(32799999),
		Claimed: false,
	},
	{
		Address: "elys1zurca5reqn4ncxf2fgej7extjsjl8rxw4jv50j",
		Amount:  math.NewInt(32799999),
		Claimed: false,
	},
	{
		Address: "elys18zj2wxaqhpvgnnh2ysx68unqwx6m4h22fyf3ps",
		Amount:  math.NewInt(32700000),
		Claimed: false,
	},
	{
		Address: "elys1q0dpqj7n4adn8f088ud2edwzeqdewar564mr9n",
		Amount:  math.NewInt(32700000),
		Claimed: false,
	},
	{
		Address: "elys1w7m2ey3lurazqc6ctz7084ysd07fm0pqrqy8ne",
		Amount:  math.NewInt(32700000),
		Claimed: false,
	},
	{
		Address: "elys1wd6kgz0nlvh3mtcehl6rqcc364lv06wattk27u",
		Amount:  math.NewInt(32700000),
		Claimed: false,
	},
	{
		Address: "elys12m47a3dyap54dyalurcn7492ntesaapmrudezt",
		Amount:  math.NewInt(32600000),
		Claimed: false,
	},
	{
		Address: "elys1swf3l8e4kqrde2mlrr8y9a3yec5etrssjr4uca",
		Amount:  math.NewInt(32500000),
		Claimed: false,
	},
	{
		Address: "elys1zkxjll43lu0ay7hezznlt0yc9uk37srrcfg6g3",
		Amount:  math.NewInt(32400000),
		Claimed: false,
	},
	{
		Address: "elys153p50xncx5zy9ne46z6w5xhmeemqw30pq57h6x",
		Amount:  math.NewInt(32299999),
		Claimed: false,
	},
	{
		Address: "elys17ah83sygpjr3dwypjwfevfmjaxvkrpmdaelxx2",
		Amount:  math.NewInt(32299999),
		Claimed: false,
	},
	{
		Address: "elys19ue0gasteqkg7nw9ty4wkyeer4g5kn38qntrrj",
		Amount:  math.NewInt(32299999),
		Claimed: false,
	},
	{
		Address: "elys1dna9t8vsh9x2t6sgjqs04236lnwq95k3rgk3x8",
		Amount:  math.NewInt(32200000),
		Claimed: false,
	},
	{
		Address: "elys1mhu7glzuverja2vu3k3es8adu84yq0lvkrsvn0",
		Amount:  math.NewInt(32200000),
		Claimed: false,
	},
	{
		Address: "elys177dduk75q2le4glpstdv0dv740va9mjshk3rs5",
		Amount:  math.NewInt(32100000),
		Claimed: false,
	},
	{
		Address: "elys1cd6kq86x6r0zt8d8a9rvaa246t9pnks4emnfgj",
		Amount:  math.NewInt(32100000),
		Claimed: false,
	},
	{
		Address: "elys1d504gcrctdg5x853ytrmkzg5wdwmqctjuhag57",
		Amount:  math.NewInt(32100000),
		Claimed: false,
	},
	{
		Address: "elys1k6jqwfnny6xvs2tva2vmt09naprpxew38vpxzz",
		Amount:  math.NewInt(32100000),
		Claimed: false,
	},
	{
		Address: "elys10qhzcmzz2uf948gw6q995y9tqfs8dape4qm7cw",
		Amount:  math.NewInt(32000000),
		Claimed: false,
	},
	{
		Address: "elys1ga8tq2jv9cks6scpmrjfhst55dqvtmhqt7c09p",
		Amount:  math.NewInt(31900000),
		Claimed: false,
	},
	{
		Address: "elys1j332lwazacscxtrdf3m0lu3kv496hwqrzazn98",
		Amount:  math.NewInt(31900000),
		Claimed: false,
	},
	{
		Address: "elys1alkraz2sy6tdd5dr7z5t9czcvx2kv9js6sg2d0",
		Amount:  math.NewInt(31600000),
		Claimed: false,
	},
	{
		Address: "elys1c9mrz6vvcvyf3kacyulcxklsnfndhs8f7m07nk",
		Amount:  math.NewInt(31600000),
		Claimed: false,
	},
	{
		Address: "elys1pgj6vp3jlex3tzh0uq79sjax6t8ajcwu63xjvm",
		Amount:  math.NewInt(31600000),
		Claimed: false,
	},
	{
		Address: "elys1vyznru7tqycz8pxx8ugxgt4l6jus55ysmnh8qf",
		Amount:  math.NewInt(31600000),
		Claimed: false,
	},
	{
		Address: "elys12t7zsz2g3dstqkgsrwm0wyw8r5pmtt36dmf4ux",
		Amount:  math.NewInt(31500000),
		Claimed: false,
	},
	{
		Address: "elys120d0gj5gf9f9ayze0yy332qemvjjferaymqjt6",
		Amount:  math.NewInt(31400000),
		Claimed: false,
	},
	{
		Address: "elys14h45qvwzkydu0vc6v5n8z8uruq6d7703jxn7yn",
		Amount:  math.NewInt(31300000),
		Claimed: false,
	},
	{
		Address: "elys17mde58mmj4gsy3gycj28ajkr6ajfa5tw729whp",
		Amount:  math.NewInt(31200000),
		Claimed: false,
	},
	{
		Address: "elys1k7rlq624pf5hmyfp79kzwp67y0ljymd2ahqj7u",
		Amount:  math.NewInt(31100000),
		Claimed: false,
	},
	{
		Address: "elys1q33v0psxnsqahx89p2gzyx3djjy6vdr4ftltws",
		Amount:  math.NewInt(31100000),
		Claimed: false,
	},
	{
		Address: "elys1j4nsrws8dt0naw2vjg5cwk6cqavkh66u2gscsz",
		Amount:  math.NewInt(31000000),
		Claimed: false,
	},
	{
		Address: "elys1knsv7ztj26fu95hpddhwz37kdrk4flhyemlxvq",
		Amount:  math.NewInt(31000000),
		Claimed: false,
	},
	{
		Address: "elys1wjmhp6j9gdun6d3xhvjr5gy4h0y0h7lfhadp95",
		Amount:  math.NewInt(31000000),
		Claimed: false,
	},
	{
		Address: "elys1nsempd9ljljetpecuyc8m6agh0up8e43x7pj4v",
		Amount:  math.NewInt(30900000),
		Claimed: false,
	},
	{
		Address: "elys1nu4w7khacz6ej5rtdxd95xh7kqzpzep6cntnl0",
		Amount:  math.NewInt(30900000),
		Claimed: false,
	},
	{
		Address: "elys155n5ezmxjagjhlns9078y7uccvqfenprxw8ut6",
		Amount:  math.NewInt(30700000),
		Claimed: false,
	},
	{
		Address: "elys1dpjytnkg7dv03ujarzfkrt7fgg7j3wf7rv2d3r",
		Amount:  math.NewInt(30700000),
		Claimed: false,
	},
	{
		Address: "elys1pmgx6a69qugmtj505ncwax05k90p69f0x4ncq0",
		Amount:  math.NewInt(30700000),
		Claimed: false,
	},
	{
		Address: "elys1yxz5ymnuv0jvjyu4elkl7fup37yq6u6n6gxrjt",
		Amount:  math.NewInt(30700000),
		Claimed: false,
	},
	{
		Address: "elys155s0dv02r7hwxe9xfjm5zgdy0hrnkq4x5pkr6w",
		Amount:  math.NewInt(30600000),
		Claimed: false,
	},
	{
		Address: "elys16c8rg8fe3r42x5kgjxqj6zzj22ynm4dpnwe0sj",
		Amount:  math.NewInt(30600000),
		Claimed: false,
	},
	{
		Address: "elys16d8eqqcqf2y35fytzl3h0jh6cdsyrx0enjchmp",
		Amount:  math.NewInt(30600000),
		Claimed: false,
	},
	{
		Address: "elys1knw3wkm76tnf4lac0xcetx8htwfge89wjwpv85",
		Amount:  math.NewInt(30600000),
		Claimed: false,
	},
	{
		Address: "elys12rvnags2yqs6g6atkcwxnjpgh07sl8znat4m35",
		Amount:  math.NewInt(30400000),
		Claimed: false,
	},
	{
		Address: "elys13wmdvjrh9tz8wdc8x8hfk5lat4k2ydgwwpgh2w",
		Amount:  math.NewInt(30400000),
		Claimed: false,
	},
	{
		Address: "elys189hdk70qx55cpq3dwwy2zt7wxq74xe2x3xpvgu",
		Amount:  math.NewInt(30400000),
		Claimed: false,
	},
	{
		Address: "elys1l4zclvy4gd480v86m59v4pv3qnzn72vwt3axwy",
		Amount:  math.NewInt(30400000),
		Claimed: false,
	},
	{
		Address: "elys1yg9adamzqvfxvkcw8c8my3yurnhcjyu8rjecjq",
		Amount:  math.NewInt(30400000),
		Claimed: false,
	},
	{
		Address: "elys1zp5u5v6lkzjvgdnkulk48d4kxhnftgauygrte4",
		Amount:  math.NewInt(30300000),
		Claimed: false,
	},
	{
		Address: "elys1nfkdjxcd52wt2tvv057twth627rrmcjycu7edk",
		Amount:  math.NewInt(30200000),
		Claimed: false,
	},
	{
		Address: "elys1rucyx9fccca205txqh8rq3xkx5yckgzua6z69w",
		Amount:  math.NewInt(30200000),
		Claimed: false,
	},
	{
		Address: "elys1crvdjr9yhs6kazl5672rdces9z45uq4hfklzr2",
		Amount:  math.NewInt(30100000),
		Claimed: false,
	},
	{
		Address: "elys13ltsfc96gfd7hy44hg4l94yk07up58mjckmtqh",
		Amount:  math.NewInt(30000000),
		Claimed: false,
	},
	{
		Address: "elys14vwyxtrt72l3ha2084rghkyke3ump3rez3maku",
		Amount:  math.NewInt(30000000),
		Claimed: false,
	},
	{
		Address: "elys1vlspucfk3kj045xrvumyns6j52ju4v0qqjn4je",
		Amount:  math.NewInt(30000000),
		Claimed: false,
	},
	{
		Address: "elys17ql5yf725samvv2t22als2rvry6wpypcxwy2zj",
		Amount:  math.NewInt(29900000),
		Claimed: false,
	},
	{
		Address: "elys1ut869g0yg7mmhjht0r8tpkkch9j74x3c7h5jfz",
		Amount:  math.NewInt(29800000),
		Claimed: false,
	},
	{
		Address: "elys1yx4dl635ky7lrd902gp2megh5fr8d7wh5elxsa",
		Amount:  math.NewInt(29800000),
		Claimed: false,
	},
	{
		Address: "elys10w9n26xq2f520pagevd6u007qpuvjt3h684jj6",
		Amount:  math.NewInt(29700000),
		Claimed: false,
	},
	{
		Address: "elys19lfg677d7r2y0m0y6c62m6swvymzfyxgfr686c",
		Amount:  math.NewInt(29700000),
		Claimed: false,
	},
	{
		Address: "elys1h3jymcp38wx3yglck4qyctn7qyuhjtex8nwyc0",
		Amount:  math.NewInt(29700000),
		Claimed: false,
	},
	{
		Address: "elys1m2e9hgllm6r884qysw6ulvrw7x94v2tu8ujkzr",
		Amount:  math.NewInt(29700000),
		Claimed: false,
	},
	{
		Address: "elys1hvqm20r20savhtynpagjr7uf2ct7khra9p40qa",
		Amount:  math.NewInt(29600000),
		Claimed: false,
	},
	{
		Address: "elys1q6wp80pwl7psggj3hulvnat9m3qf56g8qmpwp6",
		Amount:  math.NewInt(29600000),
		Claimed: false,
	},
	{
		Address: "elys15p976t3a2jwedf2pfwmyjqtje7zta4yk57n7jj",
		Amount:  math.NewInt(29500000),
		Claimed: false,
	},
	{
		Address: "elys16lkuzfh3t2j376hmt5826da9r0khkkztk5xugc",
		Amount:  math.NewInt(29500000),
		Claimed: false,
	},
	{
		Address: "elys1tdtkc4y0kvr0gyfccsvhhufqv65sf2e40mljep",
		Amount:  math.NewInt(29500000),
		Claimed: false,
	},
	{
		Address: "elys10nggwjmc5j54raxzvrz4njx5gpjz249lyhk8g4",
		Amount:  math.NewInt(29400000),
		Claimed: false,
	},
	{
		Address: "elys14pwzglqzxnpyn892vam3s4hnwk9zyf2vpl0nmk",
		Amount:  math.NewInt(29400000),
		Claimed: false,
	},
	{
		Address: "elys1pglutr70jpax68yaykn6rkzmqqeemv58ly6nuw",
		Amount:  math.NewInt(29400000),
		Claimed: false,
	},
	{
		Address: "elys1qnwjxyxnkqgdpr5fphyhefzxuh266gcjzfyhsu",
		Amount:  math.NewInt(29400000),
		Claimed: false,
	},
	{
		Address: "elys1rvh9hvcvrr2qwd4fqtle452x4h28q79puh7xyz",
		Amount:  math.NewInt(29400000),
		Claimed: false,
	},
	{
		Address: "elys1vvz3qacx3g59vu483tyy7vck2z6t8lhxw7yz3u",
		Amount:  math.NewInt(29400000),
		Claimed: false,
	},
	{
		Address: "elys18xpljsv5hhwx9kskpd0tdgtpl3mpuclmjyu9qd",
		Amount:  math.NewInt(29100000),
		Claimed: false,
	},
	{
		Address: "elys1wxnh58whlv3gsyjq7fwcdhd8mzgwguj04ymmnq",
		Amount:  math.NewInt(29100000),
		Claimed: false,
	},
	{
		Address: "elys136hp58t42y2wjlysnsuy08g3mt7wumqpwvyznu",
		Amount:  math.NewInt(29000000),
		Claimed: false,
	},
	{
		Address: "elys17rcvnga5p8xlf9d0k2mpfuchp34n37zdxanmep",
		Amount:  math.NewInt(29000000),
		Claimed: false,
	},
	{
		Address: "elys1fwju88qvyrfmnyyc4utt4vl3e7mgawercyg7n4",
		Amount:  math.NewInt(29000000),
		Claimed: false,
	},
	{
		Address: "elys1sqfyv3lgjsk747tsl0mvnpks73a6xnkhh5659q",
		Amount:  math.NewInt(29000000),
		Claimed: false,
	},
	{
		Address: "elys102gzjnynzuejawsvndrnwl8tqzszyh0att7vur",
		Amount:  math.NewInt(28900000),
		Claimed: false,
	},
	{
		Address: "elys1c9704rldz48h9jhlrj9p5dfthkql7q5xud0yth",
		Amount:  math.NewInt(28900000),
		Claimed: false,
	},
	{
		Address: "elys1fklswfel5z07x9t5gu7wzprw430z3yswsr4su0",
		Amount:  math.NewInt(28900000),
		Claimed: false,
	},
	{
		Address: "elys103verdnv4pzztty50ua7gs5gj25l2gwuef2cv6",
		Amount:  math.NewInt(28800000),
		Claimed: false,
	},
	{
		Address: "elys1ekxfzyedfm34mglp48rekn97v7k8kz8v7rqynm",
		Amount:  math.NewInt(28800000),
		Claimed: false,
	},
	{
		Address: "elys1gctysa67efud575ch9yg5r3eufd5s0ghctnrr5",
		Amount:  math.NewInt(28800000),
		Claimed: false,
	},
	{
		Address: "elys1jhhdt4pnmz0fdf7yz5a5w8zy44flw0x8a8xvc6",
		Amount:  math.NewInt(28700000),
		Claimed: false,
	},
	{
		Address: "elys1uknk5r37r6nnmmgzjdxxermcjujctzvg53wqvp",
		Amount:  math.NewInt(28700000),
		Claimed: false,
	},
	{
		Address: "elys185jhy09yc62nk6lwgjc3j22udz5wztdsx7zyry",
		Amount:  math.NewInt(28600000),
		Claimed: false,
	},
	{
		Address: "elys1a6g4llt7xdc4ja7r5ndm45qa95legcdanz65s2",
		Amount:  math.NewInt(28600000),
		Claimed: false,
	},
	{
		Address: "elys1yh8df7p3wlltqxfv4k7q49sgx84sa90hgfrf4k",
		Amount:  math.NewInt(28600000),
		Claimed: false,
	},
	{
		Address: "elys1d02cfx5sr6kuktu4senvhm565arqknysslahhl",
		Amount:  math.NewInt(28500000),
		Claimed: false,
	},
	{
		Address: "elys1f28jr42kyzy3vvmac2wxnu7gpp9p0xhfzze7mc",
		Amount:  math.NewInt(28500000),
		Claimed: false,
	},
	{
		Address: "elys1n6x2fsfzugk4zz7m4r5gumauvx237x50ut9gdt",
		Amount:  math.NewInt(28500000),
		Claimed: false,
	},
	{
		Address: "elys1c7cc99snqaeqx2fnhds2xw32nj9wyfv059j8jv",
		Amount:  math.NewInt(28400000),
		Claimed: false,
	},
	{
		Address: "elys1gr4xqrhmn44a7uu4ddhdulyrngxs0w2ngvys6z",
		Amount:  math.NewInt(28400000),
		Claimed: false,
	},
	{
		Address: "elys19r4nynn9vy4pv7020t69hnzn7chymcj3ygrscd",
		Amount:  math.NewInt(28200000),
		Claimed: false,
	},
	{
		Address: "elys1dt3p47d48ys20s4wtgxr834s2rdwu3pahc8v95",
		Amount:  math.NewInt(28200000),
		Claimed: false,
	},
	{
		Address: "elys1zsw6hn5d5cjwz7js4jjm0qd33zarhrmqxtawgj",
		Amount:  math.NewInt(28200000),
		Claimed: false,
	},
	{
		Address: "elys1dvesxq06cel2x8a9e98rq94hy7d6jlu307mxlt",
		Amount:  math.NewInt(28100000),
		Claimed: false,
	},
	{
		Address: "elys10t57vsz6fwxkpd5nyh9v05xgkarqz5fjfpszhs",
		Amount:  math.NewInt(28000000),
		Claimed: false,
	},
	{
		Address: "elys16ktapmny5yu04xwfjmvqzw0tehwumst0esrn6d",
		Amount:  math.NewInt(27900000),
		Claimed: false,
	},
	{
		Address: "elys1a290lgs7n3fhx0aspcuxdmj4lq07n39vra4e29",
		Amount:  math.NewInt(27900000),
		Claimed: false,
	},
	{
		Address: "elys1z85huk8m9m6uuqxduc9ju3ekymzamqdr7gg67d",
		Amount:  math.NewInt(27900000),
		Claimed: false,
	},
	{
		Address: "elys1sxqtdzr2wuefprez9fk2grtvru7h3zz4pm6e4n",
		Amount:  math.NewInt(27800000),
		Claimed: false,
	},
	{
		Address: "elys1vz2fcvuarywcp7cej7rx3pt8v5gaz5r9z7lgks",
		Amount:  math.NewInt(27800000),
		Claimed: false,
	},
	{
		Address: "elys1zhae93zee8mhvwuy98dlmh8prphnda6ulplmx9",
		Amount:  math.NewInt(27800000),
		Claimed: false,
	},
	{
		Address: "elys1emrugzu6zle90vshjq84ds3c076au9dzchz86x",
		Amount:  math.NewInt(27700000),
		Claimed: false,
	},
	{
		Address: "elys1yuze82stcdk3md09xpsf07mx6efx2ds0jlk30j",
		Amount:  math.NewInt(27700000),
		Claimed: false,
	},
	{
		Address: "elys14ampn8lfkzqpuf8xk5detkty7j2fhl4844szsq",
		Amount:  math.NewInt(27600000),
		Claimed: false,
	},
	{
		Address: "elys1fg20vekqd5yzry6elc035rhrqx2ve4zje620e8",
		Amount:  math.NewInt(27600000),
		Claimed: false,
	},
	{
		Address: "elys1qf7lhf2t99xanjwtm0u87hyl5k3udm4z0rhar3",
		Amount:  math.NewInt(27600000),
		Claimed: false,
	},
	{
		Address: "elys1rfzkc3863h8agtgwtmdpl447tvhpljv93p8f2c",
		Amount:  math.NewInt(27600000),
		Claimed: false,
	},
	{
		Address: "elys15qnf6turgmunhkvxrx0skvgw3apcwdp8fhfhpk",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys15rlej3p2r02ng9cl8jm4f9d6cl2j7v8xap83x0",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1cpnzuhecd6zwesmwhujrejlnjw0axqmlzlyfcs",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1e28k6zv4n78w9v89hra868zskpmqngd6vstac8",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1j6d7wcz8w7v0npj6dcsjyu0tc9x7zy05794ply",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1k8g0vlfmctyqtwahrxhudksz7rgrm6ns0acr43",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1ltmwsgr992x6s3wg7m4tg8jpd4708egc0jgru8",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1rh2v05l5uc0axzdczkgkl2g7cf3g9sd4upwtry",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1s4fhr089zcrup0hkqlqj6u5ahxrtpgwl936d4j",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1v6kdn485jkrnmpmmr2cl546a5dmlvlggmw7d53",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1vmwsy6mk5anx2qd9dncv3hfk4vr0fxawtmhr8z",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1x4mhxddz68aerxwegws9heaamj2fl0uxcklsgy",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1y5zryu2m8mp4sz3a5qy7dxx5cqxpzyalxacfsq",
		Amount:  math.NewInt(27500000),
		Claimed: false,
	},
	{
		Address: "elys1pqk89l5m5s3fxks09vfv6ljnynrfcmlw8qhe4a",
		Amount:  math.NewInt(27400000),
		Claimed: false,
	},
	{
		Address: "elys1qj9j8fnudx2a0lv9ddt82x0fykm5ajk0chfput",
		Amount:  math.NewInt(27400000),
		Claimed: false,
	},
	{
		Address: "elys1u35pnkuje88srdwhd0e8gm7w239yucd2la8756",
		Amount:  math.NewInt(27400000),
		Claimed: false,
	},
	{
		Address: "elys1u69lkdk525npweej7q2tnuhhcf4w4caqrz2zv2",
		Amount:  math.NewInt(27300000),
		Claimed: false,
	},
	{
		Address: "elys1ysuun4d6uq0ptyvg7wnyue5vzdpvspamvrz8fw",
		Amount:  math.NewInt(27300000),
		Claimed: false,
	},
	{
		Address: "elys1aexea3zcwfskjuldnyh2sxggynp5qfuxhg24ac",
		Amount:  math.NewInt(27200000),
		Claimed: false,
	},
	{
		Address: "elys1ekuxj67tnfrkmxezyprr007dfc5se2hckseyh7",
		Amount:  math.NewInt(27200000),
		Claimed: false,
	},
	{
		Address: "elys1grlfe5w5yt4l3pzyt7d4e2zcznewxhw3uamkv0",
		Amount:  math.NewInt(27200000),
		Claimed: false,
	},
	{
		Address: "elys1h9kk80r23yzfjdspmlpczga3qgm5l6qs3udl2p",
		Amount:  math.NewInt(27100000),
		Claimed: false,
	},
	{
		Address: "elys1hsd9fg32h7w224fxyj2vztu30x0q9suh9lxmdy",
		Amount:  math.NewInt(27100000),
		Claimed: false,
	},
	{
		Address: "elys1l8y8eptfuzvxh52z9299r3mnm22m4kletkw606",
		Amount:  math.NewInt(27100000),
		Claimed: false,
	},
	{
		Address: "elys1rmzlf3ajjut4zr9wv6preeg2lc785gl8646z8u",
		Amount:  math.NewInt(27100000),
		Claimed: false,
	},
	{
		Address: "elys14wkc5f7fw7wamnpn0fw9h0qe3ewpzenmz68p39",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys15j8vy7dc6vh5mc8p8wkqnq7nk7tmmwgusyyndv",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys17hpceja8hfyry0me6fxlatny7gz88ha0uhvfce",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys17nd00c9pmj53vmef3clseaxm8aa2q4ll3hw8e5",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1dxy5vd2xe3s48hlwrpy9xeepqzn3jg9h94ll9h",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1gg7ynnk8g23m5nw23024xtr2gdz8e3vpddcxrs",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1ksy2rng49fw6lt68uevc3les2murz76ezzd2nh",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1leylvdm6r25fevazc2lnyehlx5v7u6g98fede8",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1pcuwqml0m29xq5wvj9kxv3qnt8nh20ndrumup8",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1pg6sn4f8gwrhfuymvtqpex3rn4cfkkf0gsn7ed",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1q0ypalse0ww8kzj24zw2puvhpa4f3sjjlefw44",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1q4mqsfgu384elxxvgvp42pm8stap5ckrftjeqa",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1q5v7rj0d42rj2q3w7w2jvjy9gc4w30fs8lkpv5",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1qms65m0hkwkchmxw97drwhqaw5m6gg9dezz4rg",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1qwp5a2un6wqgl6ynhm2783kxhcvvjca4798ey4",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1sc7jmhrll9l9vq3zucqwrhrfv89ra3qqdelk95",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1sxmg90rxchxzreq5e244nvjjj9j6nktfl2zt76",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1t76wg9uq00pgrns3nw66mrshaylxu8evqpwpsn",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1tdec7lds4822lekyhtwnzfl7d0yxc5xpfpm29f",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1tl8wt2ffzwus0r2mf9z5ku2kuzrc68xscwlthz",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1ukmeng3uwvjsktjluchsu8et754lrl7zcqwx4u",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1uks005w8ct0dg90zdcz08ymsxfsm5q2qsv5qfh",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1uyrwawgl2r09c80ews7mt7dqy6hf8592nak7er",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1w5r2dav6llug0pt0c7jv7auupz2hwpc6qlcp6r",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1zarzc80q3rmxrmrx0kd3vltg8rv7rdtx6ulfcf",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys1zvf2yj8rpetj6ysqfp2wnl5ep9u3s4zz4jqma7",
		Amount:  math.NewInt(27000000),
		Claimed: false,
	},
	{
		Address: "elys12dq9cytc55u24r2raw2puuh8xpcrtvd35274lt",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys12mt55lfsyhkdjnel47n8f4pvdd39cr5n6up0g4",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys147c9h6r6vasgkhnxc84kpcgts2lndkt2s0njvu",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys15843yeq60zzuech6y6n7k4ggmmrq7g7z5m8e5r",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys16vrznqt603cxs2ncs9mlw4c3yctgljvqe9u9yr",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys18m0hg0eeltmq2remuwd4rjrynfjaryugzs4d2e",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys18v546fnmdqap270cn4kl3kcjuguul95u6c6cvp",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys19quc5w6e23ps0k3y2c4j53997n97uh7el9te86",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1c7mxu29k4exwarrk23qd824602c0pkfmh0jj83",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1cmrxayg26fraupnx93mwl70t9w73uddwme4qmd",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1cvnpedu0dy0zm06v32qjflt8nlm4q4x3qv4awk",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1ehwltv30yjsm747wtuhw5ltdkjmactu8dweljr",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1hhef32qjjhcy7zqx8u6e4znk598v30e87yj5wl",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1ku8s5n82p7tmjrefn7yne4shcs9z5l57ldnncl",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1mupcmj7gzh5la68lzjqgq4mhs3kujgavzly7nu",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1nc4glu5ar479n4e4jjqc9zy20fdv36h5mdx2cv",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1phmw86gl702fh885ralkh6s4j6ctl79y2gp5rd",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1tcf3ykk04j37j7g3mnykp2kpd7w0een0nhh0st",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1thzfz54jusecujarzmcphtfg80ha3azjt7lnv3",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1ty4lr6jqskxva23wgdths5ltwcxh2z4lwfmdj7",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1uxzpx4rth83kwa47wpvr04r4khl9v5w4ygn6q2",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1y56764tn3zca67dd0228y60c453f0aunmpdhal",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1ypprku64th2r0rlq93qkftjj39mlw2za4l9h8y",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1zjypjurz8eajx9zmhjzj03kgmffmd68w0l8cvh",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1zs7c6u4w39e2tp3hya3pdguhwr7547ma3vhk6t",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys1zupv0s36ewvrsw3emlvz043xrru8wjflhx87ct",
		Amount:  math.NewInt(26900000),
		Claimed: false,
	},
	{
		Address: "elys15vyfsdw7a5n4f0rwpcsyeuyv92qwlee46v8kz9",
		Amount:  math.NewInt(26800000),
		Claimed: false,
	},
	{
		Address: "elys179qsc9u03tvun4jud6swhfvv4vhq63d0lkf6m3",
		Amount:  math.NewInt(26800000),
		Claimed: false,
	},
	{
		Address: "elys1j6g8fatd588k068zfw3ayarngdghrxyywe6k64",
		Amount:  math.NewInt(26800000),
		Claimed: false,
	},
	{
		Address: "elys1k0h6cyprwua2unsfrtpdjz20x5dw6cr89y6ck3",
		Amount:  math.NewInt(26800000),
		Claimed: false,
	},
	{
		Address: "elys1lcr2dcdcc64cmhylkh7ua6wfuxjuc7l0pap5sx",
		Amount:  math.NewInt(26800000),
		Claimed: false,
	},
	{
		Address: "elys1m6qg3s3g9wpfur59dzq44e7g9hk39j9et23dgj",
		Amount:  math.NewInt(26800000),
		Claimed: false,
	},
	{
		Address: "elys1md4al6mzvpctr3053szp2rasmx0hp6rhspvnyd",
		Amount:  math.NewInt(26800000),
		Claimed: false,
	},
	{
		Address: "elys15pj342wz3kll448tfhh8l3y8fejm3xepmnu8w4",
		Amount:  math.NewInt(26700000),
		Claimed: false,
	},
	{
		Address: "elys1674t9havleu33zxcq2udsc458t3hpxtuqr889w",
		Amount:  math.NewInt(26700000),
		Claimed: false,
	},
	{
		Address: "elys18tsqvk9c4m3dpk7mu879chhk2c9z57rwgxunju",
		Amount:  math.NewInt(26700000),
		Claimed: false,
	},
	{
		Address: "elys1h3q277zx63u3mgc76a2c5t8eqzj9cc2uwj9fz9",
		Amount:  math.NewInt(26700000),
		Claimed: false,
	},
	{
		Address: "elys1lcqfxma5peq9s5kcxxvdn2lx0la28vqyy38pe2",
		Amount:  math.NewInt(26700000),
		Claimed: false,
	},
	{
		Address: "elys1lq7s8y54kv3svrqzcatw32dsvzmuzgfg4nmeg2",
		Amount:  math.NewInt(26700000),
		Claimed: false,
	},
	{
		Address: "elys1mv29uezmpgew95wkawaqvduz3p2q7svtkgsdaa",
		Amount:  math.NewInt(26700000),
		Claimed: false,
	},
	{
		Address: "elys1pjt5wnuw652qd8zee9jm7ahae6rs8d3ptv5zj4",
		Amount:  math.NewInt(26700000),
		Claimed: false,
	},
	{
		Address: "elys1rffd6vv83mjnerg802vgtgx8d5vuyqh6gn7ppl",
		Amount:  math.NewInt(26700000),
		Claimed: false,
	},
	{
		Address: "elys1yw3whsdnwwqd048trm5hkan2789mqzqtkxed6x",
		Amount:  math.NewInt(26700000),
		Claimed: false,
	},
	{
		Address: "elys1zyhpk4uka5uymulaf4eep57evca4cypff0st7r",
		Amount:  math.NewInt(26700000),
		Claimed: false,
	},
	{
		Address: "elys174wqryd5j9alllvstp8csfe4azw88hcweuyzle",
		Amount:  math.NewInt(26600000),
		Claimed: false,
	},
	{
		Address: "elys1l75nt7ypqqfycdxlzdr8h907crj2paslulasmg",
		Amount:  math.NewInt(26600000),
		Claimed: false,
	},
	{
		Address: "elys1lwrggumzeftefdj8slacayh3ppqln5uxjaqhl3",
		Amount:  math.NewInt(26600000),
		Claimed: false,
	},
	{
		Address: "elys1t2gl4lkevx4jxm4p6g5d3vc0h0j872jggvy7kk",
		Amount:  math.NewInt(26600000),
		Claimed: false,
	},
	{
		Address: "elys10tw2wdmdfzffssl6l9c4nmc98wqmumgelazttu",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys12jn3da04qwvdned7gq7awl3yre5t9f9dzh5uzv",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys13yt5g00ryqzg7wgrc72l5003eqapwvtksfw9g7",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys14zdka2942f6sr6ees04vlqwsyez4w4y99xz5qk",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys15k9q7amprt05kvphqqu7w84v94609tx32r4cxp",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys18xu279fd8xsqphmv977cmuyfwy2jl7zfrpqens",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1dyk535g9vl2ggt0xd96kpr0wypzq9fa0vnu6vk",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1kjumvvzhctwtp8pk8kdeayjurwpy3l4cgkx0hx",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1lwv7w92myhmxv70z3cywlwz7r7fpak6ejdgp0d",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1m6cph259g82adm98m0t20g58ymfug696egvdtk",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1m806kg5093nqhuxnfkjpu7v7g7pcr7yaul6x79",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1m89tgljs30k3smrkkwy7gua9e4psdjavt52wrr",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1nfz6suughv78ctn58t8h7d2wuevqq883k5g0e4",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1tw3t52v7u40tuy5wv4e76j9fmhmkt66fspqd8n",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1vjjjktsmmys3u6x9cmfjfpcpsfhkg8808frq8m",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1wx5keq3zlj2t0k0vk3r35gm7u8x8cam8etvqcm",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1yyxnneajf3fyjajjntyvc9matccg5kxh8vhjq6",
		Amount:  math.NewInt(26500000),
		Claimed: false,
	},
	{
		Address: "elys1fmla3604gd6z74xwm922zshnjkq7f09k0xdtgy",
		Amount:  math.NewInt(26400000),
		Claimed: false,
	},
	{
		Address: "elys1nys609q5wkh6t3ewa0nzqmh4ne4u5cck4vgtdr",
		Amount:  math.NewInt(26400000),
		Claimed: false,
	},
	{
		Address: "elys1tfx5eytfh9nqcjsttv7f8dwry57hakss3v7uqz",
		Amount:  math.NewInt(26400000),
		Claimed: false,
	},
	{
		Address: "elys17ew4zk3cpnpr4fsz0acc7unzxu877hek7fzrgs",
		Amount:  math.NewInt(26300000),
		Claimed: false,
	},
	{
		Address: "elys1j0uw6an5epmkx5drudyrl7j5pq3h9ucq3d86u9",
		Amount:  math.NewInt(26300000),
		Claimed: false,
	},
	{
		Address: "elys1vfczjv4u4nc3fznt4t0f3gpj9nlhy64k7r2f27",
		Amount:  math.NewInt(26300000),
		Claimed: false,
	},
	{
		Address: "elys1hvxjndz3j7jfla6cux8l5xegwj9kfcgc6ek64e",
		Amount:  math.NewInt(26200000),
		Claimed: false,
	},
	{
		Address: "elys1xerhq69xejlwf8p8jpfjdyhrnfs2ssuz5yyy2j",
		Amount:  math.NewInt(26200000),
		Claimed: false,
	},
	{
		Address: "elys1zvlah58pxpe2lpavvke9ejmzukgh5dtgz9eapf",
		Amount:  math.NewInt(26200000),
		Claimed: false,
	},
	{
		Address: "elys125hkvggsp3hh6na3js7fzzujs2k523gwwlqlw7",
		Amount:  math.NewInt(26000000),
		Claimed: false,
	},
	{
		Address: "elys1e84e9gkftuwx8pg9w4qax7f70nhf20gtq2nr3d",
		Amount:  math.NewInt(26000000),
		Claimed: false,
	},
	{
		Address: "elys1mud00tqp4u99wsc7rcvwp40n9mxgrwzm9ez0aq",
		Amount:  math.NewInt(26000000),
		Claimed: false,
	},
	{
		Address: "elys10kglmyau0df6tsktawhk8yqzp44hzvwhmh28zd",
		Amount:  math.NewInt(25900000),
		Claimed: false,
	},
	{
		Address: "elys1ae36shdnd4ny09gmx4lq9dwjsej33ajysmxq8f",
		Amount:  math.NewInt(25900000),
		Claimed: false,
	},
	{
		Address: "elys1qq6thwayczcpehdy38tr3a3t2lzkzuu4mta6v2",
		Amount:  math.NewInt(25900000),
		Claimed: false,
	},
	{
		Address: "elys1tph3d7awcqa6qxua8amzv6v66w7zksaj8y9jz9",
		Amount:  math.NewInt(25800000),
		Claimed: false,
	},
	{
		Address: "elys1u6e8mmz8wd7vcywlem6wa9w754elvw53s804px",
		Amount:  math.NewInt(25800000),
		Claimed: false,
	},
	{
		Address: "elys1x6cx2h0rqutgnt3tqy544k4ud7vjrkrfn3f8dx",
		Amount:  math.NewInt(25800000),
		Claimed: false,
	},
	{
		Address: "elys1ys4v5cktlj37wjla0qdmh2dhga8udfs3xjyk2u",
		Amount:  math.NewInt(25800000),
		Claimed: false,
	},
	{
		Address: "elys1pqk8rq4jx4kztuumr4422ehw6wvhzf6s87qmnx",
		Amount:  math.NewInt(25700000),
		Claimed: false,
	},
	{
		Address: "elys1q3h824tgm20ltynj2v29kt5sdxx34laxnfzgtt",
		Amount:  math.NewInt(25700000),
		Claimed: false,
	},
	{
		Address: "elys1qd9x93cpyqhxty0xvf4nf3cdtrsn3s3gcwxqj7",
		Amount:  math.NewInt(25700000),
		Claimed: false,
	},
	{
		Address: "elys1u29nrt6z6ad2vrvzupw06c9put8fwkart8mk23",
		Amount:  math.NewInt(25700000),
		Claimed: false,
	},
	{
		Address: "elys1w7r5zqqf4kpm8ezlzh7j6564v56nd6wugcwlp2",
		Amount:  math.NewInt(25700000),
		Claimed: false,
	},
	{
		Address: "elys1whkq75n5cycw5982u6ts39snck66e3uwqyrqgt",
		Amount:  math.NewInt(25700000),
		Claimed: false,
	},
	{
		Address: "elys1ws9n38jg3vzvwwkzytpgheymcv86srz4f3z09k",
		Amount:  math.NewInt(25700000),
		Claimed: false,
	},
	{
		Address: "elys1ncqvkdrvu8g903lc6kk3chzxdmuqkk8adt0uvs",
		Amount:  math.NewInt(25600000),
		Claimed: false,
	},
	{
		Address: "elys10we5ht6x5z9z53xdanpj9p4f3vuj0x803zce9d",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys125axgfm30wxg8jlza9ru7pu4hy3tvq2h4vr7rg",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys12jzgk2u6ku5uzpyl3flpzyegs4ewkmvh9q8666",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys12q7d5cmqznul8yanmyldgwqclgyuy5em2xwcut",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys13zcj6s9dgxjumeadpyx9wesh664x7jpe6en65s",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys15gg7tlychynfwm89lzjexe7r84j26xjxsckhp7",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys163z8u58cr82p5nm0a9m2zr482csrglxg2h0ruh",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys16cq4w94yfrfwer3xs2e7u3z37zcwd84h8k3jex",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys17s74q27389ky5ptc7xj8kh5h7x7w5x7hn6u4c2",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys1dpfa25ntxm5d8jzls6rk6vff6efft2v5qw873a",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys1eyy303757ft4ggl0wqv96krtnrszc65nphdmyg",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys1jel6m9kq57fdkzj2dd2nxlda23ygwxjg3pfkqx",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys1lf83stx6kyycdj458cq6nrc76c0lfd0kunyysw",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys1llhv3k5x0r5mlnzg9kwv63mqdtcp2vqk4fedtk",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys1n3qs3mzcjh37lp2wrgl5v865nfpue6zwv8njxh",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys1qq9ghzmx2sws4m3pnm0m5kthhp7z8wh49243z3",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys1ra5gx36gr70xmma66wcaazupz2j0jkk7lkq0p0",
		Amount:  math.NewInt(25500000),
		Claimed: false,
	},
	{
		Address: "elys17j7jx0ee4mrq7dguxmlz4cq6yhky9svfh8e9k6",
		Amount:  math.NewInt(25400000),
		Claimed: false,
	},
	{
		Address: "elys1dl689gm65gpuaxfdxn90wwhpk9xmj0kqsv5ygh",
		Amount:  math.NewInt(25400000),
		Claimed: false,
	},
	{
		Address: "elys1vcn7rt3wuc480h0xe5kp6s47my49ne2f9e28r6",
		Amount:  math.NewInt(25400000),
		Claimed: false,
	},
	{
		Address: "elys1yketqs440mm2nw4md9ngk7usxupv0gyhs0mlrr",
		Amount:  math.NewInt(25400000),
		Claimed: false,
	},
	{
		Address: "elys152wclttmddnyh978qj2fkh8efa709xu3m664es",
		Amount:  math.NewInt(25300000),
		Claimed: false,
	},
	{
		Address: "elys19dylgdkf4gqnmv95yxw024gldsg4aahy0445k7",
		Amount:  math.NewInt(25300000),
		Claimed: false,
	},
	{
		Address: "elys1jm2avxzs2jkdp0ek6kn82fgdtljfh23q2nde3d",
		Amount:  math.NewInt(25300000),
		Claimed: false,
	},
	{
		Address: "elys1889lzk75h44u2n379svypkxpg9g53h2ytfw9gf",
		Amount:  math.NewInt(25200000),
		Claimed: false,
	},
	{
		Address: "elys1drsx36tmnqzj2u9nreyrfzwsc5fyx2z9ueftzq",
		Amount:  math.NewInt(25200000),
		Claimed: false,
	},
	{
		Address: "elys1jg9j5eglm9d9ehj8ztdud3hh7zcafqalynlsr2",
		Amount:  math.NewInt(25200000),
		Claimed: false,
	},
	{
		Address: "elys1c6fh7hz786z5euwrtawxrjwjupfup92f4mwse0",
		Amount:  math.NewInt(25100000),
		Claimed: false,
	},
	{
		Address: "elys1r6knf70m2dlfljh9aqk87ldwvxshajj3muvn0p",
		Amount:  math.NewInt(25100000),
		Claimed: false,
	},
	{
		Address: "elys17tsgz565pvq72twgn8xvjfhq9ev9ctys4v280j",
		Amount:  math.NewInt(25000000),
		Claimed: false,
	},
	{
		Address: "elys17v8w550v4rxgftsey4aqtfjjla0xf7ccyfylqx",
		Amount:  math.NewInt(25000000),
		Claimed: false,
	},
	{
		Address: "elys1faq6hlver9u77rl5fdxcuke8lugl9537nk0yrx",
		Amount:  math.NewInt(25000000),
		Claimed: false,
	},
	{
		Address: "elys1pvmcefzegayc2k7msdxzzz7ekqlqn9yql3xeyc",
		Amount:  math.NewInt(25000000),
		Claimed: false,
	},
	{
		Address: "elys1x3tpfk3x9ea8qns94u8r4cp2n0jjalk30xg9a7",
		Amount:  math.NewInt(25000000),
		Claimed: false,
	},
	{
		Address: "elys17kqzpezkxw4q7fpnxgde52ggdedfzr2ydyv7ph",
		Amount:  math.NewInt(24900000),
		Claimed: false,
	},
	{
		Address: "elys13vqdazr9yqyyv3ufntxpmjedw8yda0wznhumxz",
		Amount:  math.NewInt(24800000),
		Claimed: false,
	},
	{
		Address: "elys189ekly0mmvkuq87wr05zpcj7493jupcfyk4ukj",
		Amount:  math.NewInt(24700000),
		Claimed: false,
	},
	{
		Address: "elys1ysk827mx43fdjxkxn3m9wlhn889urdyt85w3u2",
		Amount:  math.NewInt(24700000),
		Claimed: false,
	},
	{
		Address: "elys1jvqjkplredcxu7frdml7rptec0p7t3pjf4j2a2",
		Amount:  math.NewInt(24600000),
		Claimed: false,
	},
	{
		Address: "elys1032u3zt4ts8qvzd6wujn88aefs7vd0um2x9rxe",
		Amount:  math.NewInt(24500000),
		Claimed: false,
	},
	{
		Address: "elys10mew5jsqy9zr833qpvsdtmjhn59lu9vrc4v2vw",
		Amount:  math.NewInt(24500000),
		Claimed: false,
	},
	{
		Address: "elys13drf63gxsm9wh4404596vha3re2yw7v6gay3pe",
		Amount:  math.NewInt(24500000),
		Claimed: false,
	},
	{
		Address: "elys15ukppk7kzwur35xgp56jplvrnudd0fns4g8l44",
		Amount:  math.NewInt(24500000),
		Claimed: false,
	},
	{
		Address: "elys16fhx36f29yx0qr2cta9vrdjcgcts4387yk8hxe",
		Amount:  math.NewInt(24500000),
		Claimed: false,
	},
	{
		Address: "elys19jcsqajnzasalwr5ne67jefg5uszjv6jhr5n7k",
		Amount:  math.NewInt(24500000),
		Claimed: false,
	},
	{
		Address: "elys1au80cxyvz2asufwd7jhk4u9nunpxuw4l479dgk",
		Amount:  math.NewInt(24500000),
		Claimed: false,
	},
	{
		Address: "elys1l6cpln2m5yrx057k6pg2gk3ecvredd9rmwew3h",
		Amount:  math.NewInt(24500000),
		Claimed: false,
	},
	{
		Address: "elys1sp2h4vhmz80j6axjpkg0uqv6p0xtuavge32wa3",
		Amount:  math.NewInt(24500000),
		Claimed: false,
	},
	{
		Address: "elys1vlex8y05d4c0zl2g5fzwj2wggvwhnurraly0xj",
		Amount:  math.NewInt(24500000),
		Claimed: false,
	},
	{
		Address: "elys1hwwng8le73ts6hndz9j32rdax7jt3l6c2lfy42",
		Amount:  math.NewInt(24400000),
		Claimed: false,
	},
	{
		Address: "elys1yckmt2rgmvs3fhrz24z7ysphm9kg4ycm5rrdj9",
		Amount:  math.NewInt(24400000),
		Claimed: false,
	},
	{
		Address: "elys108y5j4sfg9x8hngy0wcm6c2zu4stzhttdadjc6",
		Amount:  math.NewInt(24300000),
		Claimed: false,
	},
	{
		Address: "elys109847876ejm8fuplpp9hzs8a32wsyvtjnvyck7",
		Amount:  math.NewInt(24300000),
		Claimed: false,
	},
	{
		Address: "elys1dm87s0tdqfmjntddqx3zmek547w2pk6jx9ue2m",
		Amount:  math.NewInt(24300000),
		Claimed: false,
	},
	{
		Address: "elys1eljyvvy7gqy6qe4ck98762zxh9nlfzsezczftw",
		Amount:  math.NewInt(24200000),
		Claimed: false,
	},
	{
		Address: "elys1xvrdqdvk8fke0nprcpl4qsa57c6mkk007ytqp9",
		Amount:  math.NewInt(24200000),
		Claimed: false,
	},
	{
		Address: "elys16ns4prnzmcq06ee7afudswp8yjaua7mcmkf6lz",
		Amount:  math.NewInt(24100000),
		Claimed: false,
	},
	{
		Address: "elys1mnc2cn3l4nwm632mavzdxeccwyzm879kcel4t2",
		Amount:  math.NewInt(24100000),
		Claimed: false,
	},
	{
		Address: "elys1sm3enf37nskdeswrtwlkyl22dajkmd02wg7fqp",
		Amount:  math.NewInt(24100000),
		Claimed: false,
	},
	{
		Address: "elys133vnfnmccug3lvn2s9fs2mpr56eq2hl9knntzf",
		Amount:  math.NewInt(24000000),
		Claimed: false,
	},
	{
		Address: "elys1cjwzcv7a79m4sd5642h9pgqlnyltyjyagfvvq0",
		Amount:  math.NewInt(24000000),
		Claimed: false,
	},
	{
		Address: "elys1cxful6l52eh0aknkufwtwmre7c4lvhv3fy35hh",
		Amount:  math.NewInt(24000000),
		Claimed: false,
	},
	{
		Address: "elys1s2cgwzxwyw5jkedz2rekvetnrmwkpc6cr53lc6",
		Amount:  math.NewInt(24000000),
		Claimed: false,
	},
	{
		Address: "elys1vy2f4vsmtd9w0pf8lnjqvpxwaeryxn7qrzcvrg",
		Amount:  math.NewInt(24000000),
		Claimed: false,
	},
	{
		Address: "elys1zzsn0lm456luq0h4ea654yalm52v7wthw8lxss",
		Amount:  math.NewInt(24000000),
		Claimed: false,
	},
	{
		Address: "elys18xsl630vwl5lmexqrh2p24wvrwmvqvhan8l39r",
		Amount:  math.NewInt(23900000),
		Claimed: false,
	},
	{
		Address: "elys144c9t2x0wzm6h7mfvjgdv7egdzawlhrhvk7e92",
		Amount:  math.NewInt(23800000),
		Claimed: false,
	},
	{
		Address: "elys16gqwax69966prpuzwlj5xm34prs6tcyzzq2nqc",
		Amount:  math.NewInt(23800000),
		Claimed: false,
	},
	{
		Address: "elys17e0gk06rz0cuacnm5dp4ya5te25a3mcwh0vj84",
		Amount:  math.NewInt(23800000),
		Claimed: false,
	},
	{
		Address: "elys1mts858hutnkqxrv835kdt9x5rchk8jg7mr3cpp",
		Amount:  math.NewInt(23800000),
		Claimed: false,
	},
	{
		Address: "elys1rzsft59nuuqdmhhvrr8rzq5p376mv6p5ltp04e",
		Amount:  math.NewInt(23800000),
		Claimed: false,
	},
	{
		Address: "elys1tg36hggw9l8zvv0sdedn046p743qke32ycas23",
		Amount:  math.NewInt(23800000),
		Claimed: false,
	},
	{
		Address: "elys16aw6p3spwuqlmhhmeexfnj30y9wytrl3g43750",
		Amount:  math.NewInt(23700000),
		Claimed: false,
	},
	{
		Address: "elys1cv6nktwhwdgevqr7tp85s55pe3sc93uen4waqf",
		Amount:  math.NewInt(23700000),
		Claimed: false,
	},
	{
		Address: "elys1ft84ntenrs4wr8s5qqrzxy863hxd3mzp5se3w6",
		Amount:  math.NewInt(23600000),
		Claimed: false,
	},
	{
		Address: "elys1gc55cr8hyeu7wffwq0fw4ta0xpf39w0xln83n7",
		Amount:  math.NewInt(23600000),
		Claimed: false,
	},
	{
		Address: "elys1l8durt2t2egqhpju94v7tur88xfg35lpc8pd9h",
		Amount:  math.NewInt(23600000),
		Claimed: false,
	},
	{
		Address: "elys1n545xa0qd4mhsk0hxz33l2ccctq5v5ty309s2w",
		Amount:  math.NewInt(23600000),
		Claimed: false,
	},
	{
		Address: "elys13frjl9q5hmwx2nnxgrf4v5tapv8vy53x4366gz",
		Amount:  math.NewInt(23500000),
		Claimed: false,
	},
	{
		Address: "elys1dny6gdnpnfh28z83sll4le3eegyk90lzr22eht",
		Amount:  math.NewInt(23500000),
		Claimed: false,
	},
	{
		Address: "elys1m6sge2y4swvt45sapew7wy89zlfmsqtm478k0x",
		Amount:  math.NewInt(23500000),
		Claimed: false,
	},
	{
		Address: "elys102y0s0xajewjdwwa9ghxdx73kvlkhs5fvv0pln",
		Amount:  math.NewInt(23400000),
		Claimed: false,
	},
	{
		Address: "elys129nl5t2377whqdc970jx9893kucqtkz5qygrgf",
		Amount:  math.NewInt(23400000),
		Claimed: false,
	},
	{
		Address: "elys12l7rtnz52ap803gvahftquz0s2023lpms9x9nf",
		Amount:  math.NewInt(23400000),
		Claimed: false,
	},
	{
		Address: "elys1rvngvscf9egzxuuf9amzxxwza45j6a50fganxf",
		Amount:  math.NewInt(23400000),
		Claimed: false,
	},
	{
		Address: "elys176g00la5km8s00ukewdd4adm5st9c9yln7rhsj",
		Amount:  math.NewInt(23300000),
		Claimed: false,
	},
	{
		Address: "elys1a66akx4m4c767tcewxx2hxum5jf5lazcga7dfd",
		Amount:  math.NewInt(23300000),
		Claimed: false,
	},
	{
		Address: "elys1c6852nlzkwvhtv9xr0rtchrumzgxsru0m99m7g",
		Amount:  math.NewInt(23300000),
		Claimed: false,
	},
	{
		Address: "elys1ey6p6d3urju8s7f4n4twla7yq24pp3xu3a9ry5",
		Amount:  math.NewInt(23300000),
		Claimed: false,
	},
	{
		Address: "elys1h8zwgtz3mscuhwjs87r900puresusqh9ul7uqf",
		Amount:  math.NewInt(23300000),
		Claimed: false,
	},
	{
		Address: "elys1krfm3hgv8tfjq7kqfpscqasw9tyrqr3z8shdt4",
		Amount:  math.NewInt(23300000),
		Claimed: false,
	},
	{
		Address: "elys1vrnaueftv7hewld0jftgx3nzm9en95te2yj7m0",
		Amount:  math.NewInt(23300000),
		Claimed: false,
	},
	{
		Address: "elys152nmkryu4sg7ed32vju8dggqwe2e5szel0ctfk",
		Amount:  math.NewInt(23200000),
		Claimed: false,
	},
	{
		Address: "elys1e7fepm35w0gskwxsf2awsjhy2ca579tqjf4nwn",
		Amount:  math.NewInt(23200000),
		Claimed: false,
	},
	{
		Address: "elys1njxc7wd5j6jf2653n344r9v7877ew82pfyvwgx",
		Amount:  math.NewInt(23200000),
		Claimed: false,
	},
	{
		Address: "elys1smcsdercufsyl602fys4lg6jm4q4qvxqmzvsnj",
		Amount:  math.NewInt(23200000),
		Claimed: false,
	},
	{
		Address: "elys1v7j3j46myc38h57wduafqfwsna09f3f7lj8a69",
		Amount:  math.NewInt(23200000),
		Claimed: false,
	},
	{
		Address: "elys1zzzgw04v5t2f924vp4szlrzhn4v0fngs5n4v23",
		Amount:  math.NewInt(23200000),
		Claimed: false,
	},
	{
		Address: "elys122qgav8vle7sdn8hvnfk300a7rpfujclyvw4jv",
		Amount:  math.NewInt(23100000),
		Claimed: false,
	},
	{
		Address: "elys18g7jfkn3ftmy79f7fy4qmxrg56t9m0v3lq8v8d",
		Amount:  math.NewInt(23100000),
		Claimed: false,
	},
	{
		Address: "elys1jj944umh45axcy0e0hhlaxgm7glyapn9eufzh0",
		Amount:  math.NewInt(23100000),
		Claimed: false,
	},
	{
		Address: "elys1jxxg480k4peljlggh3nnf22q02ksflzv97ykz0",
		Amount:  math.NewInt(23100000),
		Claimed: false,
	},
	{
		Address: "elys1e0z89hqh0pf6tyetnnchxz09ju5dt22ldj0avq",
		Amount:  math.NewInt(23000000),
		Claimed: false,
	},
	{
		Address: "elys1skccewp4vp9s5a6s8wj4cpmxrge4p3v4d4jlv7",
		Amount:  math.NewInt(23000000),
		Claimed: false,
	},
	{
		Address: "elys1tdjhtvet6387mp7wcrg2uu0w7xvszwuvtld56z",
		Amount:  math.NewInt(23000000),
		Claimed: false,
	},
	{
		Address: "elys14x9q836u9sn6mqlrspxd46v7tscaam7sprrhry",
		Amount:  math.NewInt(22900000),
		Claimed: false,
	},
	{
		Address: "elys19l6haqxagtud5qsrhhqeqv2y6rvr6jfkq3w9wr",
		Amount:  math.NewInt(22900000),
		Claimed: false,
	},
	{
		Address: "elys19nl8t0vekf9plvppvtrtafzy46hcrh3ds05hlp",
		Amount:  math.NewInt(22900000),
		Claimed: false,
	},
	{
		Address: "elys1dah2u2k4grq7cmzhj6zn4uqdw0jueytp8yuks9",
		Amount:  math.NewInt(22900000),
		Claimed: false,
	},
	{
		Address: "elys1mv99semd6mnlgjgxwdzfvspdplkdwucxu50dgw",
		Amount:  math.NewInt(22900000),
		Claimed: false,
	},
	{
		Address: "elys1ygp5jrqxx35wwkk85amm5zc2jgyzlwts90h2sp",
		Amount:  math.NewInt(22900000),
		Claimed: false,
	},
	{
		Address: "elys163d02qdctvzj4fj0qr2788als2laau2cs38t3z",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys17s883k8ygp4ldea6yj7yq6tjv78vzwkplfes7t",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys19e3lcuv452a99pthwapm56uqgj4jwc87qz3tgn",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys19gvwtwwx2g6hnua7smrkvdhvz04qxunrmx64qx",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys19zrgkct6pfjux4sfl3wgkfpawl9nk0c0lmks8f",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys1duygsz025w69pnyvyvfc2dddce6a23v7c07kvk",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys1frnvgksrmnyuz8ufwmv8j639dtg5fu4xtvmgaw",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys1kt67k0x7wggf0qgdszdzj2ycf3yp4rjhldttq6",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys1kuwl702jatawz4p46p3324f9wcy340x3v6x2g4",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys1sn6hyvqr7y2fdsm6w4gldq76ez55xk3fc785gr",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys1u6v4fqcvqrgma0y04frt9xaz229jyx55fk24j5",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys1vemmttg4c0lsqwsdhrjc24ktaahfnh9h8ru9t9",
		Amount:  math.NewInt(22800000),
		Claimed: false,
	},
	{
		Address: "elys13kfud2fgkah9recll2q045jeuugj7jng0ykysc",
		Amount:  math.NewInt(22700000),
		Claimed: false,
	},
	{
		Address: "elys1auz63ya078c6rvy4heya9auzk069qn9c2w7y3q",
		Amount:  math.NewInt(22700000),
		Claimed: false,
	},
	{
		Address: "elys1cxxdetuy0e7sjws6z2993xjkuk8yeqyfjj4eh7",
		Amount:  math.NewInt(22700000),
		Claimed: false,
	},
	{
		Address: "elys1rcty3se5m8ll66wtfuxsgekaf2xsgd2lggt5c7",
		Amount:  math.NewInt(22700000),
		Claimed: false,
	},
	{
		Address: "elys1swtgsydywt3f4fyjn7y2f09ffaz8y6xs56uf3s",
		Amount:  math.NewInt(22700000),
		Claimed: false,
	},
	{
		Address: "elys1x93xgkggdwdm5un8ryf646ccnpl6jnhj7f2ynv",
		Amount:  math.NewInt(22700000),
		Claimed: false,
	},
	{
		Address: "elys1x94dd6n7f286ay239cxzkpua2gyqfcq996saqt",
		Amount:  math.NewInt(22700000),
		Claimed: false,
	},
	{
		Address: "elys100d70t6pqcgt3xywe9u5u73gxwe3zzqp7rkyye",
		Amount:  math.NewInt(22600000),
		Claimed: false,
	},
	{
		Address: "elys14rp0w0xn4q6luxml8u82zuf7twkmv7ussunfv4",
		Amount:  math.NewInt(22600000),
		Claimed: false,
	},
	{
		Address: "elys1ndp7vk03p72p0uwtupx9cynxw2ma4v8mwqa09j",
		Amount:  math.NewInt(22600000),
		Claimed: false,
	},
	{
		Address: "elys1tthnejnsc4t57m40dg3xd3shtl6et5z8fuqfvn",
		Amount:  math.NewInt(22600000),
		Claimed: false,
	},
	{
		Address: "elys19xz7d8wjpwxc8xfp8y55qvdsqpsnpql0fedtkn",
		Amount:  math.NewInt(22500000),
		Claimed: false,
	},
	{
		Address: "elys1wmyrc30j7u4vyk02rd33wvpm83e8yt2ekakfed",
		Amount:  math.NewInt(22500000),
		Claimed: false,
	},
	{
		Address: "elys15prkc0l5v2384fmlh43hs4k9krwj2zx03sl6d3",
		Amount:  math.NewInt(22400000),
		Claimed: false,
	},
	{
		Address: "elys15rwsu3yyj0ch0ewd9qfegcuf9melxp8jrns6fk",
		Amount:  math.NewInt(22400000),
		Claimed: false,
	},
	{
		Address: "elys1e0fdmtn8h94kz6k3qhzch575k78axwwlx4z9xk",
		Amount:  math.NewInt(22400000),
		Claimed: false,
	},
	{
		Address: "elys1kwpr7y9y95qdjmauu8gu5mwmjgnwtwmly2cmzs",
		Amount:  math.NewInt(22400000),
		Claimed: false,
	},
	{
		Address: "elys1lk0dc9tulnvklsa486dqrhk59wk3l02vs5cagc",
		Amount:  math.NewInt(22400000),
		Claimed: false,
	},
	{
		Address: "elys1s6lg2hyusjkkw40gfx59f3jkzyc7w7pv4j48ff",
		Amount:  math.NewInt(22400000),
		Claimed: false,
	},
	{
		Address: "elys1ty2r0zfrr8ngdyzqf5hhmy7zs96udh5fuqt8c2",
		Amount:  math.NewInt(22400000),
		Claimed: false,
	},
	{
		Address: "elys1uqz09jduerkccr3fyqltwtpsmpu0uspvxmlxry",
		Amount:  math.NewInt(22400000),
		Claimed: false,
	},
	{
		Address: "elys1w32h0fevk2r4j6pl2a8yxagsczdw097reduwrq",
		Amount:  math.NewInt(22400000),
		Claimed: false,
	},
	{
		Address: "elys1yyt4mre8ajjrcxudy3r9fq4ue6jefhxstqyngn",
		Amount:  math.NewInt(22400000),
		Claimed: false,
	},
	{
		Address: "elys1zsdugrgggkpzkpsxxs8vz6hfynujnej3gf4hwx",
		Amount:  math.NewInt(22400000),
		Claimed: false,
	},
	{
		Address: "elys14q2ae3dr0arslnvjlnds9vqz04mxhjmaruet37",
		Amount:  math.NewInt(22300000),
		Claimed: false,
	},
	{
		Address: "elys16hfcr3xfr0c0jt3lmqg5qfczrtltvx6shrqc67",
		Amount:  math.NewInt(22300000),
		Claimed: false,
	},
	{
		Address: "elys16sjjfaf2chhv839kufmanratv09kn3rp6rxqjl",
		Amount:  math.NewInt(22300000),
		Claimed: false,
	},
	{
		Address: "elys1na04c6lftr9ljp4uslq8me8w3kc6uwa3rv0406",
		Amount:  math.NewInt(22300000),
		Claimed: false,
	},
	{
		Address: "elys13ftqq7rzumakddkqhc4ze28ndly434ld5qxrfs",
		Amount:  math.NewInt(22200000),
		Claimed: false,
	},
	{
		Address: "elys16evzag2kdkl94t7zjmdmj78045p8zehc6c0z2f",
		Amount:  math.NewInt(22200000),
		Claimed: false,
	},
	{
		Address: "elys1g3ndee4ru6d3wrdvz54zpz5f7t9vdyecdk6wvu",
		Amount:  math.NewInt(22200000),
		Claimed: false,
	},
	{
		Address: "elys1kg0mrkap8jp9ahff4dxnnzhfus3amj5vtqm29y",
		Amount:  math.NewInt(22200000),
		Claimed: false,
	},
	{
		Address: "elys106q82p0xm9ayrctzef3jjrgr5q8wuy6rvnmk4d",
		Amount:  math.NewInt(22100000),
		Claimed: false,
	},
	{
		Address: "elys1a338296u9htmlkwk98x87fkxxm5jrt0k8mfppd",
		Amount:  math.NewInt(22100000),
		Claimed: false,
	},
	{
		Address: "elys1fa786pzjcx4guuzheezmal7c0suv7dluqte0e5",
		Amount:  math.NewInt(22100000),
		Claimed: false,
	},
	{
		Address: "elys1mt7m5cg44kgenrazu5g5zxx3g79zsq080hr89s",
		Amount:  math.NewInt(22100000),
		Claimed: false,
	},
	{
		Address: "elys1cm5kpysu05vxzdp60nr55uaxplmwkvqu6jntcx",
		Amount:  math.NewInt(22000000),
		Claimed: false,
	},
	{
		Address: "elys1fns3m2sr8wlhfrsn75p05n586fqc44nxn8wd9c",
		Amount:  math.NewInt(22000000),
		Claimed: false,
	},
	{
		Address: "elys1fvhrv7f4jp04zjw3a65jeaxeq5tlx2dfy70xmq",
		Amount:  math.NewInt(22000000),
		Claimed: false,
	},
	{
		Address: "elys16ztsl9mxyx30jchyzse0z70evc56w9xyv8wfkm",
		Amount:  math.NewInt(21900000),
		Claimed: false,
	},
	{
		Address: "elys1aguke87dzyss5d5p5pjh2uufqfapvkd0tg3r39",
		Amount:  math.NewInt(21900000),
		Claimed: false,
	},
	{
		Address: "elys1aqlx5v4reyw7k239vensz9s9ugehrwgenshqew",
		Amount:  math.NewInt(21900000),
		Claimed: false,
	},
	{
		Address: "elys1q355l6ewfh5jqx3r5ufhng7kujvldtxzk7c2ag",
		Amount:  math.NewInt(21900000),
		Claimed: false,
	},
	{
		Address: "elys1mjalqrl5rt4ga2dksy4e77zg57dw2aksnprptt",
		Amount:  math.NewInt(21800000),
		Claimed: false,
	},
	{
		Address: "elys12sztn5lmwqd2nlemmefj0uwy6e092kus5cj8qg",
		Amount:  math.NewInt(21700000),
		Claimed: false,
	},
	{
		Address: "elys1gu6dyquykt6hwx90hxctnq7v28vf3eapzys886",
		Amount:  math.NewInt(21700000),
		Claimed: false,
	},
	{
		Address: "elys1mjjyhsdpvtcnw99xcvpxwz94ry4fccjgkjge5c",
		Amount:  math.NewInt(21700000),
		Claimed: false,
	},
	{
		Address: "elys1x0awrmm9cp0trgv3wx5mkr9mdmmc8yymzx3l0j",
		Amount:  math.NewInt(21700000),
		Claimed: false,
	},
	{
		Address: "elys1zknchf92a9hn3xgvyut4m3xdymqzxt4c4c38dp",
		Amount:  math.NewInt(21700000),
		Claimed: false,
	},
	{
		Address: "elys13z746ztwxgq7dc6ylygw6g7vqz8qg36l2uuxyz",
		Amount:  math.NewInt(21600000),
		Claimed: false,
	},
	{
		Address: "elys186sj02nuxn4fk8g79kf9mafuftz34erx0r9d2l",
		Amount:  math.NewInt(21600000),
		Claimed: false,
	},
	{
		Address: "elys1hake9gt024lmpagdx50knuztcyssuqdkdf7xq9",
		Amount:  math.NewInt(21600000),
		Claimed: false,
	},
	{
		Address: "elys1kl274uw594e47m6akkppyk57rrd2p6xmjc5jks",
		Amount:  math.NewInt(21600000),
		Claimed: false,
	},
	{
		Address: "elys1uj5em8w6ljt7a05v7y36fczpm42ly5mf344l5f",
		Amount:  math.NewInt(21600000),
		Claimed: false,
	},
	{
		Address: "elys1v5vjcksvqskutlv5rv6j9a9umwg4ufk7ughghf",
		Amount:  math.NewInt(21600000),
		Claimed: false,
	},
	{
		Address: "elys10vjucyxnx77ve0qq7c244wl7lcnadazknq37up",
		Amount:  math.NewInt(21500000),
		Claimed: false,
	},
	{
		Address: "elys1392vmkuqyzy4upqlwuu42nfq004gfgw2hknjns",
		Amount:  math.NewInt(21500000),
		Claimed: false,
	},
	{
		Address: "elys1qdfrxy3kxhmev485e0ezyss324tzdzg6z25s96",
		Amount:  math.NewInt(21500000),
		Claimed: false,
	},
	{
		Address: "elys1tfj0mdn2qj5nyltm9jj2lrshvm3x3dchjax4pz",
		Amount:  math.NewInt(21500000),
		Claimed: false,
	},
	{
		Address: "elys1u99uv5yujh8ajl39kaxk5ndpax0m6l3tj82xza",
		Amount:  math.NewInt(21500000),
		Claimed: false,
	},
	{
		Address: "elys1vkqq0q26hmxkuf9xq3fdwsdxp702uffvjzwv4e",
		Amount:  math.NewInt(21500000),
		Claimed: false,
	},
	{
		Address: "elys13hqgyynx0nsl289lqhrdytg4rpad09kmdsuraw",
		Amount:  math.NewInt(21400000),
		Claimed: false,
	},
	{
		Address: "elys18u9teze2gcqz6xsna8hhvf4fj80xcx6sklzmtk",
		Amount:  math.NewInt(21400000),
		Claimed: false,
	},
	{
		Address: "elys16d9p3d29kkgf40pc04yjmkf7n4hmk78rsusjgp",
		Amount:  math.NewInt(21300000),
		Claimed: false,
	},
	{
		Address: "elys17tcjjwh6lw24h2ck3yf3rp8mv3sennvhjvz4t0",
		Amount:  math.NewInt(21300000),
		Claimed: false,
	},
	{
		Address: "elys1dg4kvp76gn6h2sz3xexqeyk66pudpfz4mpwt53",
		Amount:  math.NewInt(21300000),
		Claimed: false,
	},
	{
		Address: "elys1vvyewh6c2yx2fy0auqgckql8zqxvx8d3n3lnpz",
		Amount:  math.NewInt(21300000),
		Claimed: false,
	},
	{
		Address: "elys17gktyn3nmuzmh0n5fjmkkrcqyh55y4d7wktnks",
		Amount:  math.NewInt(21200000),
		Claimed: false,
	},
	{
		Address: "elys1c0zraqkq4nm845qwm4kgprw24fg3k5czacfrye",
		Amount:  math.NewInt(21200000),
		Claimed: false,
	},
	{
		Address: "elys1fdfhjfev5475p7vn5fp2l3g5rx0s6xpxzj5gys",
		Amount:  math.NewInt(21200000),
		Claimed: false,
	},
	{
		Address: "elys1nmu56xwlyc00vvpaalg8eaqjz4pyde8e368r3h",
		Amount:  math.NewInt(21200000),
		Claimed: false,
	},
	{
		Address: "elys16mpfg7w5c85584zkvq25e07mam7lwh9xgw3tar",
		Amount:  math.NewInt(21100000),
		Claimed: false,
	},
	{
		Address: "elys1gx8ndz8llvlz3szp2htytkjfzcxznen7ke7dal",
		Amount:  math.NewInt(21100000),
		Claimed: false,
	},
	{
		Address: "elys1tjszxjner6l9z3szzfy9rd98ky2q8hhqtkvk96",
		Amount:  math.NewInt(21000000),
		Claimed: false,
	},
	{
		Address: "elys10naft24cu8lv7ppkmpwjc9hdkaulsjerpf02l3",
		Amount:  math.NewInt(20900000),
		Claimed: false,
	},
	{
		Address: "elys166d4y46h679ygmy8cad7qkqacut5yylletrkv8",
		Amount:  math.NewInt(20900000),
		Claimed: false,
	},
	{
		Address: "elys1785plfjck87fmvglzjkggzaclc6hd5vvwddarz",
		Amount:  math.NewInt(20900000),
		Claimed: false,
	},
	{
		Address: "elys19twrud2kz3jh7vxylwqxp5r82c7xkju7eds9tr",
		Amount:  math.NewInt(20900000),
		Claimed: false,
	},
	{
		Address: "elys1ch20ap73dqudanw0ufhpn2jpvl6zayqeac99l2",
		Amount:  math.NewInt(20900000),
		Claimed: false,
	},
	{
		Address: "elys1lcfgsnhlgxtq8dklaljz87d4yak666k5uee9q7",
		Amount:  math.NewInt(20900000),
		Claimed: false,
	},
	{
		Address: "elys1ruvympw45vhnqc7upweaep5dgrvwvc2yx4qhjs",
		Amount:  math.NewInt(20900000),
		Claimed: false,
	},
	{
		Address: "elys1sc4z4qam4pyguhg8exy4g6vzp43fgm7cp4jfec",
		Amount:  math.NewInt(20900000),
		Claimed: false,
	},
	{
		Address: "elys12kquqtqu2lx9pnpycr00t3g5ucdx9krjsjkc43",
		Amount:  math.NewInt(20800000),
		Claimed: false,
	},
	{
		Address: "elys12tycshmw6vmn4mhq6m4xqzs6gfg5czgszlavea",
		Amount:  math.NewInt(20800000),
		Claimed: false,
	},
	{
		Address: "elys12uxs3pxc4vjw5hzjrfj37u6napyelfk2sg6slg",
		Amount:  math.NewInt(20800000),
		Claimed: false,
	},
	{
		Address: "elys18edf5rgrwm5wfhrmayy3n8jug34j5c845es0tn",
		Amount:  math.NewInt(20800000),
		Claimed: false,
	},
	{
		Address: "elys18qcsnl6n59fmv0df6yf6n2c8cnj0rll5jzd9qd",
		Amount:  math.NewInt(20800000),
		Claimed: false,
	},
	{
		Address: "elys1f0kcpeuys5mhdcpq87xmgh4vny8vwkl0ag32q9",
		Amount:  math.NewInt(20800000),
		Claimed: false,
	},
	{
		Address: "elys1fvwhl5ecvw0n6vcrupg32lcj57es9fwtww5wlz",
		Amount:  math.NewInt(20800000),
		Claimed: false,
	},
	{
		Address: "elys1hetx9a3q9cujy3h5x8zk43lu8krreup7d8lk26",
		Amount:  math.NewInt(20800000),
		Claimed: false,
	},
	{
		Address: "elys1ystqhs6r9pmfce5jfyntpxu8apdg6eh48546jh",
		Amount:  math.NewInt(20800000),
		Claimed: false,
	},
	{
		Address: "elys1fvpardtvufppc0mvuuttqt6tnkgq32md4dfaxr",
		Amount:  math.NewInt(20700000),
		Claimed: false,
	},
	{
		Address: "elys1nqc7cz8nv4jqcwh0ww8kvm7nv3cc5q27tstten",
		Amount:  math.NewInt(20700000),
		Claimed: false,
	},
	{
		Address: "elys1qgjp0j0ujakd2eqx8vdqnttx0pdf2mph3zxqcs",
		Amount:  math.NewInt(20700000),
		Claimed: false,
	},
	{
		Address: "elys1spwgpszc5g0056sdr78ugfqxahjcjqgdrr6cw9",
		Amount:  math.NewInt(20700000),
		Claimed: false,
	},
	{
		Address: "elys1t2syrj7s78kgswwzkyw3hnne9q9p6zmues024v",
		Amount:  math.NewInt(20700000),
		Claimed: false,
	},
	{
		Address: "elys1vxjrsywc6dkussjnam7tcgqfhm27d93zxwaxaj",
		Amount:  math.NewInt(20700000),
		Claimed: false,
	},
	{
		Address: "elys17w5mmv9yn5x38nxekphrt0hhuffva5v6dyemn3",
		Amount:  math.NewInt(20600000),
		Claimed: false,
	},
	{
		Address: "elys18q7x5pgk6q0jr8c875jxqurm4ftatfas7uuq70",
		Amount:  math.NewInt(20600000),
		Claimed: false,
	},
	{
		Address: "elys1mhaxlqdd5mxa2rjs0qggu3lrztwym89vqewvvl",
		Amount:  math.NewInt(20600000),
		Claimed: false,
	},
	{
		Address: "elys1nyw8tgh7j62y99w4t4wkn7gu977rdrqka9dz3x",
		Amount:  math.NewInt(20600000),
		Claimed: false,
	},
	{
		Address: "elys1uc0h5nzn64avkfuj2dm4glxem9xu30l980xufm",
		Amount:  math.NewInt(20600000),
		Claimed: false,
	},
	{
		Address: "elys1vjymkucluukngrzyp3ykkdrj6mr6h7y0uev9vw",
		Amount:  math.NewInt(20600000),
		Claimed: false,
	},
	{
		Address: "elys1g5tq363ned62yrt3t9e3kseq4ydzxf4mgaruf2",
		Amount:  math.NewInt(20500000),
		Claimed: false,
	},
	{
		Address: "elys1jv903newfsm6m94v3aatdvxzzkrz03fp3gumne",
		Amount:  math.NewInt(20500000),
		Claimed: false,
	},
	{
		Address: "elys1vdqhpe6w2rukld8l0l64l2yww2a36y3w7lmq2c",
		Amount:  math.NewInt(20500000),
		Claimed: false,
	},
	{
		Address: "elys10axt2ywuggtax8zqjq5f9l4at9hk228qdlj3e2",
		Amount:  math.NewInt(20400000),
		Claimed: false,
	},
	{
		Address: "elys146sfsp8pkjx8ajysc52jnr2tczg44p3228atxh",
		Amount:  math.NewInt(20400000),
		Claimed: false,
	},
	{
		Address: "elys17npt8r6pn4ztxaw3ywksyzkap8ns2e6r4j28c9",
		Amount:  math.NewInt(20400000),
		Claimed: false,
	},
	{
		Address: "elys1aqp8qvnaqtfda0n2s2f92uvmvkmg5keeyrccsh",
		Amount:  math.NewInt(20400000),
		Claimed: false,
	},
	{
		Address: "elys1fu2p20fw8fhm6cclzqal6t5jd9q4yx49mz62yh",
		Amount:  math.NewInt(20400000),
		Claimed: false,
	},
	{
		Address: "elys1g9ymv7w4gzf6xll97v9rr62dw7als0t4tp29kg",
		Amount:  math.NewInt(20400000),
		Claimed: false,
	},
	{
		Address: "elys1lt8mxtnwpdkc3yjcwmww62u673nu3744nxg0ma",
		Amount:  math.NewInt(20400000),
		Claimed: false,
	},
	{
		Address: "elys1t4ta20hs975rt5pqxklljdflqrv93gp7d8p8af",
		Amount:  math.NewInt(20400000),
		Claimed: false,
	},
	{
		Address: "elys1285wwlq7ctkevp52lj5st5xst0trq02r99pwt6",
		Amount:  math.NewInt(20300000),
		Claimed: false,
	},
	{
		Address: "elys12na9kpdskkmg5jzpqje70pv33lydkqjnh8lwyj",
		Amount:  math.NewInt(20300000),
		Claimed: false,
	},
	{
		Address: "elys1sl0mzv7vgtulq7ypjcxm9h5s9wckz2zj62u8rx",
		Amount:  math.NewInt(20300000),
		Claimed: false,
	},
	{
		Address: "elys12ku3c7ynevdsv97uj5typ5rhv92n3m4tjrsf4s",
		Amount:  math.NewInt(20100000),
		Claimed: false,
	},
	{
		Address: "elys1770v2mejrqvr5hel3j4qd5eym9rzegdqg2t8ft",
		Amount:  math.NewInt(20100000),
		Claimed: false,
	},
	{
		Address: "elys1wktj9glqjr6kvxj0m4ws340qulpzkzv7zzrj97",
		Amount:  math.NewInt(20100000),
		Claimed: false,
	},
	{
		Address: "elys14ku40qtlmjuss54pa89h43kt3v847j708hqw6y",
		Amount:  math.NewInt(20000000),
		Claimed: false,
	},
	{
		Address: "elys1gypsw5k3m5f72qpha8sdew96ernmtfz0xpskfh",
		Amount:  math.NewInt(20000000),
		Claimed: false,
	},
	{
		Address: "elys1jvhe858g3zq99qaceg0npldvzcyvc07d8puh0r",
		Amount:  math.NewInt(20000000),
		Claimed: false,
	},
	{
		Address: "elys1mvhe8e7ws6ck83fj40enk8zza2tudl8xlrpsc9",
		Amount:  math.NewInt(20000000),
		Claimed: false,
	},
	{
		Address: "elys1pprg96g7mym0tlt89vy6hgudfv4gs5rqcl8gzh",
		Amount:  math.NewInt(20000000),
		Claimed: false,
	},
	{
		Address: "elys1q9vygyflge0enpu5tr3m62mjzu7nwjlyc3gcx3",
		Amount:  math.NewInt(20000000),
		Claimed: false,
	},
	{
		Address: "elys1vuw42wdskl6haj9srn4853y9d3ptwd75mxhk0p",
		Amount:  math.NewInt(20000000),
		Claimed: false,
	},
	{
		Address: "elys1kkehkd575ugq0qrwn0ylah0e2rmzp9yfra96yf",
		Amount:  math.NewInt(19900000),
		Claimed: false,
	},
	{
		Address: "elys1mgnjvrw3cds9af7wpjks6a64vq83r2zzgr7wa7",
		Amount:  math.NewInt(19900000),
		Claimed: false,
	},
	{
		Address: "elys19z6akh568lk46pdrx7ty4p85z3fcfyhmlyljcs",
		Amount:  math.NewInt(19800000),
		Claimed: false,
	},
	{
		Address: "elys1dxd3c8hk6xcy83uqy797f4j3rdclnkd80mzayn",
		Amount:  math.NewInt(19800000),
		Claimed: false,
	},
	{
		Address: "elys1lkjw9sva8yt5a9fys2wnrudvrxgyk977myw8fe",
		Amount:  math.NewInt(19800000),
		Claimed: false,
	},
	{
		Address: "elys1nllnffudccy2f6wrfv4djaadsvqsr2lacvvntx",
		Amount:  math.NewInt(19800000),
		Claimed: false,
	},
	{
		Address: "elys1rr8l9jmcg9zjc9crm83l7r8zdkdudws6w07pvj",
		Amount:  math.NewInt(19800000),
		Claimed: false,
	},
	{
		Address: "elys1sw666up8y0yqr5vlkw0rt7zztd92jcngzmde5p",
		Amount:  math.NewInt(19800000),
		Claimed: false,
	},
	{
		Address: "elys10w6pdxd7d4p9ep4dhjf52aka7lqgfck3yr0mgh",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys12enk3m3jld54hmf72yd3l6qhczghh49wnuy33x",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys14m5k8ysjwx73zt33pzl44gxf5q0amakrkmpnmq",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys14q4c6kgg7xegcc88naqt6f8fz2skuzczrhxac6",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys154q2qf8j4gqdhr7ujac2rluez8hazlxl6w4sf7",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys16xaqgqauajujgjppctetc32fqdwp0r6sql2ekr",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys19s9870729lngap2kzls3dcyeu7a4m2tcslhmp8",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys1elw7dfkar9wgg5em7ddsz6u3t3yw9wg82mwpj6",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys1fnp5wdq4l9ed287698xcykwlj0g0dmcu6u3p66",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys1kwtrcsg5huktkwzkfwwgvecrmwzlxzs7mkgfnh",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys1nq7aems5sgfggak2lsu8nn74cwj48jrspffdzp",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys1q3tz46plr7hsx3qzdssv5f5wy94av3m3wzj80e",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys1zhv6enwzn04gsx70755xluxs6n62fwxjkph8nt",
		Amount:  math.NewInt(19700000),
		Claimed: false,
	},
	{
		Address: "elys10hta62cyl2j9jcv766ruyda3ytt2wq8kn5u58u",
		Amount:  math.NewInt(19600000),
		Claimed: false,
	},
	{
		Address: "elys18h0puzgu0jthrhq69t7r6sf7clhv7j9z9jrcth",
		Amount:  math.NewInt(19600000),
		Claimed: false,
	},
	{
		Address: "elys1jrg50jer2l0h56lxup292yyz2u0lfajtdzc9mg",
		Amount:  math.NewInt(19600000),
		Claimed: false,
	},
	{
		Address: "elys1k0ffqwunzuljt48qmsc3evgcssfzxhtl85qa8s",
		Amount:  math.NewInt(19600000),
		Claimed: false,
	},
	{
		Address: "elys1l5sl3dz9yf3w6785vuf25wa3pacm7fl4ja4ngu",
		Amount:  math.NewInt(19600000),
		Claimed: false,
	},
	{
		Address: "elys1m39qcw67tpa9vk89xmllgkc5llvd93t2qjwhpl",
		Amount:  math.NewInt(19600000),
		Claimed: false,
	},
	{
		Address: "elys1mgll8lfmpdyadctt2ga5vza0dtw3gtpyvmd68q",
		Amount:  math.NewInt(19600000),
		Claimed: false,
	},
	{
		Address: "elys1vc7h2kmlw4a36y2ssmgaavm3t0tpfgwm7c9yk3",
		Amount:  math.NewInt(19600000),
		Claimed: false,
	},
	{
		Address: "elys1wr7kfspwa5dtn7j64mqw0ea6snhcy8whj3h6pv",
		Amount:  math.NewInt(19600000),
		Claimed: false,
	},
	{
		Address: "elys1yfh5gms9zfqlqgngznzvtvz9juy6t3kvcpn7gj",
		Amount:  math.NewInt(19600000),
		Claimed: false,
	},
	{
		Address: "elys17wcmnlpz3kclez8vrfk00al8xg2up3auzvtzjv",
		Amount:  math.NewInt(19500000),
		Claimed: false,
	},
	{
		Address: "elys1drg6dztvhkga7uxft4vm0urfad3kaw2hmnl7wy",
		Amount:  math.NewInt(19500000),
		Claimed: false,
	},
	{
		Address: "elys1kcjjgf88x3m083hf8whhjrcafwy6n8zghn2u35",
		Amount:  math.NewInt(19500000),
		Claimed: false,
	},
	{
		Address: "elys1szggzp3082m9y0mm7t430ucnl2ccpcsuew8m00",
		Amount:  math.NewInt(19500000),
		Claimed: false,
	},
	{
		Address: "elys1vag5ynn027l03fd9ej80pyq4ynmg8uuzt8nqzh",
		Amount:  math.NewInt(19500000),
		Claimed: false,
	},
	{
		Address: "elys187dpsc6ar8d96unhvt2mqhv5upn3mdy6zxkh7x",
		Amount:  math.NewInt(19400000),
		Claimed: false,
	},
	{
		Address: "elys19sw8vghdjm3unh5yl9sfm3h0qz7umdqn0vhpqv",
		Amount:  math.NewInt(19400000),
		Claimed: false,
	},
	{
		Address: "elys1elpwvlv4fn43mdykvllxqe277pz3a6sqenqjcx",
		Amount:  math.NewInt(19400000),
		Claimed: false,
	},
	{
		Address: "elys1eqxn7d8uuxflz3hl9ws0w35r96k029tqdh3v4h",
		Amount:  math.NewInt(19400000),
		Claimed: false,
	},
	{
		Address: "elys1sqdak3jxuqq0xdgd94vhjwrtqazwpj5qfx56gr",
		Amount:  math.NewInt(19400000),
		Claimed: false,
	},
	{
		Address: "elys1sx08xdxefe9xzauldx9km6fvffjjz2z6k85me8",
		Amount:  math.NewInt(19400000),
		Claimed: false,
	},
	{
		Address: "elys1zud3h5akuch5360w8jgmc4qd8dcpush8mz6yk3",
		Amount:  math.NewInt(19400000),
		Claimed: false,
	},
	{
		Address: "elys12uk22nzee0hgahzttujcdce78ax627as8wp0xq",
		Amount:  math.NewInt(19300000),
		Claimed: false,
	},
	{
		Address: "elys1e7mlt0e429yrz7rgh45gytkj29cgf0kt399efd",
		Amount:  math.NewInt(19300000),
		Claimed: false,
	},
	{
		Address: "elys1xayhxtey4z6np2gfvm7c3drnxtexlalsyxz4dp",
		Amount:  math.NewInt(19300000),
		Claimed: false,
	},
	{
		Address: "elys1ls0zmrkfer27tx5grll8zhxln50uwgf5q0cvfn",
		Amount:  math.NewInt(19200000),
		Claimed: false,
	},
	{
		Address: "elys1rwsxqxrr0leg0tamysc4xfypyu2ngxkq5dzs6a",
		Amount:  math.NewInt(19200000),
		Claimed: false,
	},
	{
		Address: "elys1zz39g5yc2crgjzv83tg08wzd4hl27zqgte2kc0",
		Amount:  math.NewInt(19200000),
		Claimed: false,
	},
	{
		Address: "elys169a473gxnuujazf55faf984hmxwu9hqwnw26cr",
		Amount:  math.NewInt(19100000),
		Claimed: false,
	},
	{
		Address: "elys19e92zzazdfu99s5q2743xhcgraauejx8t5ee2e",
		Amount:  math.NewInt(19100000),
		Claimed: false,
	},
	{
		Address: "elys1m35wt8w0h9gqqrea97xrxl4xx5qyphsjxv2xu2",
		Amount:  math.NewInt(19100000),
		Claimed: false,
	},
	{
		Address: "elys1r8nplz585hux0atndjegxfnkuvwynsu4jmmxvu",
		Amount:  math.NewInt(19100000),
		Claimed: false,
	},
	{
		Address: "elys1tcszdp03jqjk3jfusgxscw3ncu6vqc0qkzdqq8",
		Amount:  math.NewInt(19100000),
		Claimed: false,
	},
	{
		Address: "elys1yrqwr8ql9j6mr7yc3fmtmzelzesjgfln0w7qwt",
		Amount:  math.NewInt(19100000),
		Claimed: false,
	},
	{
		Address: "elys10tejargmx0wggedejg7q0lqa2rtcqqdxqf6pap",
		Amount:  math.NewInt(19000000),
		Claimed: false,
	},
	{
		Address: "elys142amxx2z488pk84ys765wh2hfm3n5r6728hm96",
		Amount:  math.NewInt(19000000),
		Claimed: false,
	},
	{
		Address: "elys143ysyrmzf3qzmwst7m48wrvanz4mzgn7sd55ly",
		Amount:  math.NewInt(19000000),
		Claimed: false,
	},
	{
		Address: "elys18gwty582aujvl75erznzcv9gzvwg6pvxcytu2q",
		Amount:  math.NewInt(19000000),
		Claimed: false,
	},
	{
		Address: "elys1qwted69szn32huygszhdme7yqs5a3jaf8fj3rw",
		Amount:  math.NewInt(19000000),
		Claimed: false,
	},
	{
		Address: "elys1y888s6gvxcuszwkyxf26leympx6wuta6lkedyr",
		Amount:  math.NewInt(19000000),
		Claimed: false,
	},
	{
		Address: "elys12j6yyuxfmkp3hd7asvra8mhcv67f4v7tfqyxaj",
		Amount:  math.NewInt(18900000),
		Claimed: false,
	},
	{
		Address: "elys1c2c7v5e7f77r25239uxxfg829xs404xt0pax9l",
		Amount:  math.NewInt(18900000),
		Claimed: false,
	},
	{
		Address: "elys1fxq6uqd7vegk8hq5m74g8qnug5fh4kea9xl7rp",
		Amount:  math.NewInt(18900000),
		Claimed: false,
	},
	{
		Address: "elys1hgqe65s2x9azrkuwa659teuv3m2j0rwt365pft",
		Amount:  math.NewInt(18900000),
		Claimed: false,
	},
	{
		Address: "elys1ken7mreuwt59x3mafsxxkwfp6sjh0lspgdfjkk",
		Amount:  math.NewInt(18900000),
		Claimed: false,
	},
	{
		Address: "elys1l9c2h64gzljc8v9lcy8qqupzjeaveq0gy0xk3g",
		Amount:  math.NewInt(18900000),
		Claimed: false,
	},
	{
		Address: "elys1nwxhtsw0es28xwg5xt4c48s733tcr24vwx7n8a",
		Amount:  math.NewInt(18900000),
		Claimed: false,
	},
	{
		Address: "elys1p9tj3dj5903hm6237g8k6j3y4gfm0x8azpar05",
		Amount:  math.NewInt(18900000),
		Claimed: false,
	},
	{
		Address: "elys1t9249fxy9pzmf9vs0yhww0yv7wq2xcffmdlfxl",
		Amount:  math.NewInt(18900000),
		Claimed: false,
	},
	{
		Address: "elys18urmck7n3ks35syu40l9p45z537dsqk4glx0ly",
		Amount:  math.NewInt(18800000),
		Claimed: false,
	},
	{
		Address: "elys1rasm3y5tkvzqaf7rvhy246js02lk589xx60pv3",
		Amount:  math.NewInt(18800000),
		Claimed: false,
	},
	{
		Address: "elys1vx7wlc4npgj39f73lte395mhv9x4gfy83a8le0",
		Amount:  math.NewInt(18800000),
		Claimed: false,
	},
	{
		Address: "elys128w8qtnpchawwqraymhnvc0ygg6ecwur3pawmm",
		Amount:  math.NewInt(18700000),
		Claimed: false,
	},
	{
		Address: "elys18cpfx8x3lg5p96vwpe8ncw25a0466jrujn8d4w",
		Amount:  math.NewInt(18700000),
		Claimed: false,
	},
	{
		Address: "elys1hyp76v3fjljv0cerq59atda45e8qjhak59t7hk",
		Amount:  math.NewInt(18700000),
		Claimed: false,
	},
	{
		Address: "elys1n2sgn4j4a4fv3cyflhz5ujg52a2qcphax2qfga",
		Amount:  math.NewInt(18600000),
		Claimed: false,
	},
	{
		Address: "elys1tp3xv7xdwq4wvn53ms9k2j7e43jemvqdg55085",
		Amount:  math.NewInt(18600000),
		Claimed: false,
	},
	{
		Address: "elys12s6rp0rtsxm3vw45n4svw6s48frkm32s7x7ztz",
		Amount:  math.NewInt(18500000),
		Claimed: false,
	},
	{
		Address: "elys1uzc6a27n99r2juugl05am0e4c5f5q4zt5vn5jh",
		Amount:  math.NewInt(18500000),
		Claimed: false,
	},
	{
		Address: "elys1yntny92rqnehm4uaaxls67e2yacdlzw2gm3lzs",
		Amount:  math.NewInt(18500000),
		Claimed: false,
	},
	{
		Address: "elys14qu603gvjh3kn2rl9mtqa0fl0m29wtgdphv0cc",
		Amount:  math.NewInt(18400000),
		Claimed: false,
	},
	{
		Address: "elys1gdnnfzt2mdega5wkkp9uvfw98pqkfrllc05zj7",
		Amount:  math.NewInt(18400000),
		Claimed: false,
	},
	{
		Address: "elys1whyz7w57898rmvq734x4ey06zg559rp3c3d5vt",
		Amount:  math.NewInt(18400000),
		Claimed: false,
	},
	{
		Address: "elys1uhwp5m37vt5pjfxjtrcumtqdtf5rrtekc9mudv",
		Amount:  math.NewInt(18300000),
		Claimed: false,
	},
	{
		Address: "elys18f64dyzjddml73u35qjvg0fgdehp52at94dut5",
		Amount:  math.NewInt(18200000),
		Claimed: false,
	},
	{
		Address: "elys1adq6wslusnqnh5tmed9a7gkzv5ku3xe5nw53lh",
		Amount:  math.NewInt(18200000),
		Claimed: false,
	},
	{
		Address: "elys1esgw83ay3kmvyzmvz78329n4mgv44xwjdnqeen",
		Amount:  math.NewInt(18200000),
		Claimed: false,
	},
	{
		Address: "elys1qg9vg3vxtrvmt67e3mcwywra820ny0z7ew3450",
		Amount:  math.NewInt(18200000),
		Claimed: false,
	},
	{
		Address: "elys152edetn7dtqhzmx4ql0n8vfwfryks8ws8ysdt0",
		Amount:  math.NewInt(18100000),
		Claimed: false,
	},
	{
		Address: "elys16phjkxgruznsq8cu70ry2p50xr4fezercw4rv8",
		Amount:  math.NewInt(18100000),
		Claimed: false,
	},
	{
		Address: "elys1wjexccyqsu2vars362sgttgut30c3ldkj20txr",
		Amount:  math.NewInt(18100000),
		Claimed: false,
	},
	{
		Address: "elys1mtwspc8m3vza00j4c8v9yfazd0nhsdg79g79r0",
		Amount:  math.NewInt(18000000),
		Claimed: false,
	},
	{
		Address: "elys1u3erv2tms3c0ymagad03wsvmc7qr9l7a0jtaje",
		Amount:  math.NewInt(17900000),
		Claimed: false,
	},
	{
		Address: "elys1vq2sw5cd9hzuswsw9xc9ngpymyhnt5ens36ttn",
		Amount:  math.NewInt(17900000),
		Claimed: false,
	},
	{
		Address: "elys1qsxn3pjk9gt6fkrfq0tu68vkcxkxexnja6umem",
		Amount:  math.NewInt(17800000),
		Claimed: false,
	},
	{
		Address: "elys1k58zka6d93tzgr9wdp3t2hzaex32dum56m8n63",
		Amount:  math.NewInt(17700000),
		Claimed: false,
	},
	{
		Address: "elys1r7jummrtzex9sx98er2xyx3uj7qf3mfjehnttv",
		Amount:  math.NewInt(17700000),
		Claimed: false,
	},
	{
		Address: "elys1hqm2lxevw7en8hzhq3wwsax52982h6nv2hnsud",
		Amount:  math.NewInt(17600000),
		Claimed: false,
	},
	{
		Address: "elys1p4fllfqcx6e6w90qpmyekqulz8fcdalgqfhewt",
		Amount:  math.NewInt(17600000),
		Claimed: false,
	},
	{
		Address: "elys1u4s6jx25ezpyx7zc2tpjf2jnfktslsfpzgf7at",
		Amount:  math.NewInt(17600000),
		Claimed: false,
	},
	{
		Address: "elys10ve48hxz6f0ultm58664z0jtq80zr5wn4pwpvg",
		Amount:  math.NewInt(17500000),
		Claimed: false,
	},
	{
		Address: "elys102uur87a3xt2wep6hzpqjujjxau0ge7w89du57",
		Amount:  math.NewInt(17400000),
		Claimed: false,
	},
	{
		Address: "elys1kdgnmg7k5gfnlvw9zumteezfh6tqcpmqass4vx",
		Amount:  math.NewInt(17400000),
		Claimed: false,
	},
	{
		Address: "elys18kv2tdy7safr2cfjzh08ycysxe3q3gevr3nug7",
		Amount:  math.NewInt(17200000),
		Claimed: false,
	},
	{
		Address: "elys1asgr2udukq6x5ywwvzxea65wxc4mexwndlrhp3",
		Amount:  math.NewInt(17100000),
		Claimed: false,
	},
	{
		Address: "elys1t2xdununjt7q54ms920cukdfdzcxan9sghgycv",
		Amount:  math.NewInt(17100000),
		Claimed: false,
	},
	{
		Address: "elys1xq8trh7nwg2zntr4gqrm5gnnz4wjxc0rfthcs8",
		Amount:  math.NewInt(17000000),
		Claimed: false,
	},
	{
		Address: "elys1m9r22qp80ckvq570qd3ytthx5m8qgr326cjqgx",
		Amount:  math.NewInt(16900000),
		Claimed: false,
	},
	{
		Address: "elys1zrj5vf089y4g6xy5cgax3w3sg30832kpzchw6q",
		Amount:  math.NewInt(16900000),
		Claimed: false,
	},
	{
		Address: "elys1w78hwllnfrglyl04yaqwdtxda55jnyqze98y95",
		Amount:  math.NewInt(16800000),
		Claimed: false,
	},
	{
		Address: "elys10w8x9cxu7zl5qjpnxe5gxl0cnpqltrqvp02ac0",
		Amount:  math.NewInt(16700000),
		Claimed: false,
	},
	{
		Address: "elys14gjsvvtpng0rrzk7cjtcafy83mp3uegs4fw695",
		Amount:  math.NewInt(16700000),
		Claimed: false,
	},
	{
		Address: "elys1s4mza5j3syulgkxcz7spsy6a5r57g3gqnmgm9e",
		Amount:  math.NewInt(16700000),
		Claimed: false,
	},
	{
		Address: "elys1t94ny6lggha9m3euy560t9m6lsnck6ae48w0qk",
		Amount:  math.NewInt(16700000),
		Claimed: false,
	},
	{
		Address: "elys1ae6wgvx53tg3x5cvutsthhpt6qy6e4qv5hvmn3",
		Amount:  math.NewInt(16500000),
		Claimed: false,
	},
	{
		Address: "elys1glvnnwjmh9l4uwyjdqyl83fx4dxearqkhm25qr",
		Amount:  math.NewInt(16500000),
		Claimed: false,
	},
	{
		Address: "elys13kgy0lw7mfe0g6fsazzwlnl3vpl057h2yru5z0",
		Amount:  math.NewInt(16200000),
		Claimed: false,
	},
	{
		Address: "elys1neqkvzwujt7nhag9hnfjnvskwd3cgvpaax2g6e",
		Amount:  math.NewInt(16000000),
		Claimed: false,
	},
	{
		Address: "elys1v8adh84zq3mv6wxnqcuxwa29n0v0sga5v7uz42",
		Amount:  math.NewInt(16000000),
		Claimed: false,
	},
	{
		Address: "elys18rfen52303z55995g76lvmku5t9wk6yg06u5c6",
		Amount:  math.NewInt(15900000),
		Claimed: false,
	},
	{
		Address: "elys1kfk598ulsf9d9lem5s5zhgtngrnpqhrl6c780z",
		Amount:  math.NewInt(15900000),
		Claimed: false,
	},
	{
		Address: "elys1vqzn2kzjnhvn02x0xc5n3z2wwaw766asv4cg6q",
		Amount:  math.NewInt(15900000),
		Claimed: false,
	},
	{
		Address: "elys1wje5wl7t4u88ayffzejazeq896gj47jky76l85",
		Amount:  math.NewInt(15900000),
		Claimed: false,
	},
	{
		Address: "elys14rywxen7ywr5ang6ylu2am8ewwgnsadae7kds8",
		Amount:  math.NewInt(15800000),
		Claimed: false,
	},
	{
		Address: "elys17n2zruy7j3q9hk2cdkpnt52sw0d345d8ln3rrn",
		Amount:  math.NewInt(15800000),
		Claimed: false,
	},
	{
		Address: "elys1qs45amjamdgknnncuw9wqf3f6scvuue6ap7d6g",
		Amount:  math.NewInt(15800000),
		Claimed: false,
	},
	{
		Address: "elys1vjumuukff2lhnl6w2qyv5fma570a2ck5d0w7lz",
		Amount:  math.NewInt(15800000),
		Claimed: false,
	},
	{
		Address: "elys1e2yqsczzv9kecafzyeut6jygd54lrsvw796lj5",
		Amount:  math.NewInt(15700000),
		Claimed: false,
	},
	{
		Address: "elys1h9exkjjywdxp87ylxmynk02cfejddzcz7479uu",
		Amount:  math.NewInt(15700000),
		Claimed: false,
	},
	{
		Address: "elys1jh6pjsja4pu9s8eqra7r4c62sq3hvutcvpck3q",
		Amount:  math.NewInt(15700000),
		Claimed: false,
	},
	{
		Address: "elys1cdytdwlt2ha853x7dl05svt37wathzuz5fqyu6",
		Amount:  math.NewInt(15600000),
		Claimed: false,
	},
	{
		Address: "elys1hwmg2ksdtmlg9dwl06tg33mpmwwdc9lafyrw37",
		Amount:  math.NewInt(15600000),
		Claimed: false,
	},
	{
		Address: "elys1nazmdjdsjmjvwy8r4aqgxz2f2svs0gmnx26e33",
		Amount:  math.NewInt(15600000),
		Claimed: false,
	},
	{
		Address: "elys1stgh98ya8vaf4xe7fd42dz2tv3ya6nzj4l5yae",
		Amount:  math.NewInt(15600000),
		Claimed: false,
	},
	{
		Address: "elys18hgrcd8e3gt9r7ppm3qf7s4us26cfzdvr8arwl",
		Amount:  math.NewInt(15500000),
		Claimed: false,
	},
	{
		Address: "elys1qnzypl92yssha2hj2evne93jgnh7mz7jwk8nj8",
		Amount:  math.NewInt(15300000),
		Claimed: false,
	},
	{
		Address: "elys1fzt2hhh5whllapndqjjm3u2lq6dk0wu4f0dvl5",
		Amount:  math.NewInt(15200000),
		Claimed: false,
	},
	{
		Address: "elys1arsplxt5rqx9ce4wyns59c6w4ujwxw40yjy8rp",
		Amount:  math.NewInt(15100000),
		Claimed: false,
	},
	{
		Address: "elys1dgp7q5knvzlfw2v2ckyje89052an80pe3kgh6w",
		Amount:  math.NewInt(15100000),
		Claimed: false,
	},
	{
		Address: "elys1qxe237xvewxtezvxx9cg54qwp5p2ewhw82f8nh",
		Amount:  math.NewInt(15000000),
		Claimed: false,
	},
	{
		Address: "elys10a044pmaaz3drmtk7n4d5gcjtmnhhl9n3mv2sx",
		Amount:  math.NewInt(14900000),
		Claimed: false,
	},
	{
		Address: "elys17j2c4mvh4f65605p0nu6evfrdzt0f77yfzcg8z",
		Amount:  math.NewInt(14900000),
		Claimed: false,
	},
	{
		Address: "elys1lf6edxnn30jzec0hy8f0cg8hzzftz758xucmg9",
		Amount:  math.NewInt(14900000),
		Claimed: false,
	},
	{
		Address: "elys1prfvsxynx7w955h5heye3l52myxaqkl4f3yrcm",
		Amount:  math.NewInt(14900000),
		Claimed: false,
	},
	{
		Address: "elys1wqaluwpj8q5dfccvyuhgdsen6ld6pr42sxyx0f",
		Amount:  math.NewInt(14800000),
		Claimed: false,
	},
	{
		Address: "elys1fceqqrd9gzj7h0g6qcvz883ymuw9c4g2yxkag6",
		Amount:  math.NewInt(14600000),
		Claimed: false,
	},
	{
		Address: "elys1tpcx9s8vh0uwnt8xrg7vvjr4sd57yvsuu70jxj",
		Amount:  math.NewInt(14500000),
		Claimed: false,
	},
	{
		Address: "elys14kx7dcnfjjegtz34pmxhsrym630q7egm20wewt",
		Amount:  math.NewInt(14400000),
		Claimed: false,
	},
	{
		Address: "elys1lazt3tp3286ruh30eek563ammw4jwf8u3vned3",
		Amount:  math.NewInt(14400000),
		Claimed: false,
	},
	{
		Address: "elys1tq8wv2vzvef5d8xsxgpfyyu68hljarpggkj078",
		Amount:  math.NewInt(14400000),
		Claimed: false,
	},
	{
		Address: "elys123rtea80kdwzaywkz2wjrhacmccwne9nz54rar",
		Amount:  math.NewInt(14300000),
		Claimed: false,
	},
	{
		Address: "elys1q893nx9yehg4w5ljvuv8t7dlf9x55v0cznrlhs",
		Amount:  math.NewInt(14300000),
		Claimed: false,
	},
	{
		Address: "elys1xuh8engpzdq9q8flxj288rg35thp023rzxe944",
		Amount:  math.NewInt(14300000),
		Claimed: false,
	},
	{
		Address: "elys1h979acxultpglzcuvw7taw747cwwd9purjmzss",
		Amount:  math.NewInt(14200000),
		Claimed: false,
	},
	{
		Address: "elys1ydlaewl5gvsfe5hu0mez58er0s939tfhx9kzpe",
		Amount:  math.NewInt(14200000),
		Claimed: false,
	},
	{
		Address: "elys10sng5fj7keyp9gf8nqne3yy9strl9ugc48n550",
		Amount:  math.NewInt(14100000),
		Claimed: false,
	},
	{
		Address: "elys183xw8l4vs7jxx872hucdp56dfgansunz88y0fx",
		Amount:  math.NewInt(14100000),
		Claimed: false,
	},
	{
		Address: "elys1ha62rx3jy79w6k4dvkygtlp6guc825h8jq9vz8",
		Amount:  math.NewInt(14100000),
		Claimed: false,
	},
	{
		Address: "elys1nf3f77wnp425q8j2xq48u8at8ltv62k6h3s6t4",
		Amount:  math.NewInt(14100000),
		Claimed: false,
	},
	{
		Address: "elys1r9wtf88wxzena5680d8sftlayu02v5l8mlf3pl",
		Amount:  math.NewInt(14100000),
		Claimed: false,
	},
	{
		Address: "elys1fd5t4q4vhnntry9pmjg6djv64d58yxtah45caj",
		Amount:  math.NewInt(14000000),
		Claimed: false,
	},
	{
		Address: "elys1fypru5z563pesx7gs7q7xejxa0qswr2dessd6r",
		Amount:  math.NewInt(13900000),
		Claimed: false,
	},
	{
		Address: "elys1mk60fm5e06dh5ja06jyw0vyqdx249wc0nhcaxf",
		Amount:  math.NewInt(13900000),
		Claimed: false,
	},
	{
		Address: "elys1vkumc97gfmp4sv5vddf5aupfsczqq70ahge32v",
		Amount:  math.NewInt(13900000),
		Claimed: false,
	},
	{
		Address: "elys1w2thskvqkfqkv2wz0uvec9p2zlcrfgwl8wna36",
		Amount:  math.NewInt(13900000),
		Claimed: false,
	},
	{
		Address: "elys10zseswd0le2aasvd8n0w6m86p69klh2hz5lr96",
		Amount:  math.NewInt(13800000),
		Claimed: false,
	},
	{
		Address: "elys16lhp0r6gcw23jh4tjv8x46ynkmmk2ck4n57t0v",
		Amount:  math.NewInt(13800000),
		Claimed: false,
	},
	{
		Address: "elys1dugnyw2rkxuhvl7r744rxcqpqrdh5l8jpkce5r",
		Amount:  math.NewInt(13800000),
		Claimed: false,
	},
	{
		Address: "elys1m3xltwk5qlz0gyu554cex0vvgfyww6a7kv328r",
		Amount:  math.NewInt(13800000),
		Claimed: false,
	},
	{
		Address: "elys1x9l2zslth552mjpayey944hl6e32rzljsadeja",
		Amount:  math.NewInt(13800000),
		Claimed: false,
	},
	{
		Address: "elys1xrm5pzmlr3dpvnen6v3gysehlnq2cfh4lrz3sa",
		Amount:  math.NewInt(13800000),
		Claimed: false,
	},
	{
		Address: "elys123ypud9e0a7mlmpuy62w2ehvmz67xz64jfdeac",
		Amount:  math.NewInt(13700000),
		Claimed: false,
	},
	{
		Address: "elys16uz4kzsj8ywvmun2q3ctkh9upau9gj2rxrwvnu",
		Amount:  math.NewInt(13700000),
		Claimed: false,
	},
	{
		Address: "elys19vk6fvlppnqndcs6gezdu0arrldsaty2ejwcp8",
		Amount:  math.NewInt(13700000),
		Claimed: false,
	},
	{
		Address: "elys1ekrp4f45aur3y3eja2q2uwnnmtjwq9pdsmzc6f",
		Amount:  math.NewInt(13700000),
		Claimed: false,
	},
	{
		Address: "elys1tmhgfjmhqdg5adalulshdu7dtuhr8unrhg8pqw",
		Amount:  math.NewInt(13700000),
		Claimed: false,
	},
	{
		Address: "elys14z0kctudvjc9ksntax03k07k0q0humfagyd7ue",
		Amount:  math.NewInt(13600000),
		Claimed: false,
	},
	{
		Address: "elys1hlxurkgnwr4nrc93krlj92rxfl8hgd2ycdjrwt",
		Amount:  math.NewInt(13600000),
		Claimed: false,
	},
	{
		Address: "elys1j2wv8ewpph5dkrp4fdjt7wrutkqh02kpxvpu5q",
		Amount:  math.NewInt(13600000),
		Claimed: false,
	},
	{
		Address: "elys1qpv8fete9ln7ev4nkehj7quwdssw4v9g9hekrh",
		Amount:  math.NewInt(13600000),
		Claimed: false,
	},
	{
		Address: "elys1ure99rdg86dlh93503xyefuk0z575a8802jqjy",
		Amount:  math.NewInt(13600000),
		Claimed: false,
	},
	{
		Address: "elys1d42awqn64x5yqhsn7jksf49xxv5qqgm6znmgh5",
		Amount:  math.NewInt(13500000),
		Claimed: false,
	},
	{
		Address: "elys1msaj6sl2typmep65fhp3pmwrypn489zrl49ejy",
		Amount:  math.NewInt(13500000),
		Claimed: false,
	},
	{
		Address: "elys1vjn3lalwzxylwy4def7ukkdq7etg3nnp5q4efn",
		Amount:  math.NewInt(13500000),
		Claimed: false,
	},
	{
		Address: "elys1y5kflvms303ausemce2h3t2gkztvdfvw3rsced",
		Amount:  math.NewInt(13500000),
		Claimed: false,
	},
	{
		Address: "elys14hh9uew39es766zqp0yhsa3fephqn9qfhhftaq",
		Amount:  math.NewInt(13400000),
		Claimed: false,
	},
	{
		Address: "elys1elq32g07s9keun0p0guns06lthgj7thrdfrl4s",
		Amount:  math.NewInt(13400000),
		Claimed: false,
	},
	{
		Address: "elys1jne8akuffl79uxph8c75rgc05xef3p30mtqpy7",
		Amount:  math.NewInt(13400000),
		Claimed: false,
	},
	{
		Address: "elys1tgw6pwk828ssvkexdtwn9d3yaqvhnztvxl28wq",
		Amount:  math.NewInt(13400000),
		Claimed: false,
	},
	{
		Address: "elys1zffyupsg4ckan05w5ecdylnzx09clt55kmmppd",
		Amount:  math.NewInt(13400000),
		Claimed: false,
	},
	{
		Address: "elys1cvv0lwljj956na3txq748maflw6fyj5e90y543",
		Amount:  math.NewInt(13300000),
		Claimed: false,
	},
	{
		Address: "elys1d6ck5c4paz380vzn4nfyfkl0kdq8xy23xemkgc",
		Amount:  math.NewInt(13200000),
		Claimed: false,
	},
	{
		Address: "elys1e8d4gxtwmw7k76nk2n0tuzcdr9ud3hrk049z26",
		Amount:  math.NewInt(13200000),
		Claimed: false,
	},
	{
		Address: "elys1t8xha08ja42m2twg0yn6ct2qjzut7fudnj6k5p",
		Amount:  math.NewInt(13200000),
		Claimed: false,
	},
	{
		Address: "elys1wdux442zn9eenspg2ezru7p3xhy69skdl6kf9c",
		Amount:  math.NewInt(13200000),
		Claimed: false,
	},
	{
		Address: "elys13k7wxw39mr89rcaydc928uc9yvdtskr89exjwt",
		Amount:  math.NewInt(13100000),
		Claimed: false,
	},
	{
		Address: "elys15t08eumz43fvqgumd2rgpgkv324s94rs9rphye",
		Amount:  math.NewInt(13100000),
		Claimed: false,
	},
	{
		Address: "elys1jfudukzca4rxu9fwzcfgxjlu52m05g088vcnr9",
		Amount:  math.NewInt(13100000),
		Claimed: false,
	},
	{
		Address: "elys1tacusfldfmc97sgjjvmxzvnnhhr0h3ummg9wg0",
		Amount:  math.NewInt(13100000),
		Claimed: false,
	},
	{
		Address: "elys1zpngc2ty2mmttk060jjv80lm5r38hqjva654y5",
		Amount:  math.NewInt(13100000),
		Claimed: false,
	},
	{
		Address: "elys103qdjrglpj2u2eyfug8yznkvpxpcjz63asyu7s",
		Amount:  math.NewInt(13000000),
		Claimed: false,
	},
	{
		Address: "elys1l5favcwjnyx366rc3jqvvz4zs85795pyddug37",
		Amount:  math.NewInt(13000000),
		Claimed: false,
	},
	{
		Address: "elys1vzr4j99qhln9jgkeka9t5584nmyypyxd37t4ng",
		Amount:  math.NewInt(13000000),
		Claimed: false,
	},
	{
		Address: "elys1yd6wrgm37sr7s3xwg65uvy3fww8r73xy0v576m",
		Amount:  math.NewInt(13000000),
		Claimed: false,
	},
	{
		Address: "elys1av2ylvmdznvle4670rmsv6gyktfyt6pz34f23w",
		Amount:  math.NewInt(12800000),
		Claimed: false,
	},
	{
		Address: "elys1w53w844uyx342g5anwf0t88xhk2s8al8k3vayt",
		Amount:  math.NewInt(12800000),
		Claimed: false,
	},
	{
		Address: "elys1gfjfh62z5mgtycd3wpt9qhkw65jcykx86cevdn",
		Amount:  math.NewInt(12700000),
		Claimed: false,
	},
	{
		Address: "elys1unrdfrnef4tf90wdkhkwuadp05d43nvg6mzkk8",
		Amount:  math.NewInt(12700000),
		Claimed: false,
	},
	{
		Address: "elys13nfumc7lexuv5qd0qurlgjdf6fw2067s0fvujk",
		Amount:  math.NewInt(12600000),
		Claimed: false,
	},
	{
		Address: "elys1frmxxwasdgumawft3pktfks8r4qkk0uaq3xtwx",
		Amount:  math.NewInt(12600000),
		Claimed: false,
	},
	{
		Address: "elys1k2nwtnxrtmh4q9mjsp360lvwse52v9879s3ynn",
		Amount:  math.NewInt(12500000),
		Claimed: false,
	},
	{
		Address: "elys1u86gadg4y2m7l0t9jkpspeq7luk5kd89mjvtck",
		Amount:  math.NewInt(12500000),
		Claimed: false,
	},
	{
		Address: "elys1w9up9msmu789gqdqvr252gtw90gyvghupgnp30",
		Amount:  math.NewInt(12500000),
		Claimed: false,
	},
	{
		Address: "elys104d7pw3yurk05yyt0e5lwnexkyu2uh7tas0kud",
		Amount:  math.NewInt(12400000),
		Claimed: false,
	},
	{
		Address: "elys10gs6phgedj20dalg6qu2q2c8rgj30kf3xmwtg7",
		Amount:  math.NewInt(12300000),
		Claimed: false,
	},
	{
		Address: "elys1vdw0pv5664ml0elz45vl92vy0auyrx9vlwwtw5",
		Amount:  math.NewInt(12300000),
		Claimed: false,
	},
	{
		Address: "elys1cknjdj8hgcpqhvuk7rpsszr6zsy25fpq9w43w3",
		Amount:  math.NewInt(12200000),
		Claimed: false,
	},
	{
		Address: "elys1gs7vwm4r50rftzfymf9r0gpmkaqy0gjp4w5kqp",
		Amount:  math.NewInt(12200000),
		Claimed: false,
	},
	{
		Address: "elys1lylak4z7nzqy6eemer5twpr7ugmgehlew7g6fa",
		Amount:  math.NewInt(12200000),
		Claimed: false,
	},
	{
		Address: "elys10xj640cadu0qn9fnux7fmk29f7htx6fe30jms5",
		Amount:  math.NewInt(12100000),
		Claimed: false,
	},
	{
		Address: "elys1vtvcc3g3kn9jfnadvcgk5r3846vwc3gaqpa6rh",
		Amount:  math.NewInt(12100000),
		Claimed: false,
	},
	{
		Address: "elys1fmzv8ucpvlc20x80g85ym6x98zxw0j4gmagzkf",
		Amount:  math.NewInt(11900000),
		Claimed: false,
	},
	{
		Address: "elys1kxzz4820jn7jzwc5npddtcsafh867kc4a9cugf",
		Amount:  math.NewInt(11900000),
		Claimed: false,
	},
	{
		Address: "elys1x9gh3y47fk36jgy8p0usc2v4evf5x7fmcs4wdc",
		Amount:  math.NewInt(11900000),
		Claimed: false,
	},
	{
		Address: "elys1dekmdps6y29ps46tsx4wgw8ek67gm4wwv05etm",
		Amount:  math.NewInt(11800000),
		Claimed: false,
	},
	{
		Address: "elys1lefp4tzntgparxsx2puqffwplcsr9jg0wgysxz",
		Amount:  math.NewInt(11800000),
		Claimed: false,
	},
	{
		Address: "elys1v9w49rtdpkwxeje8xnx9pnf9z9hyr2e0vldehj",
		Amount:  math.NewInt(11700000),
		Claimed: false,
	},
	{
		Address: "elys1em7rlpdhce9pcjmrlgpukp84uvn5y8uqw87c9z",
		Amount:  math.NewInt(11600000),
		Claimed: false,
	},
	{
		Address: "elys1jsj9hewhd8dvwfckmgf632r70ajaqkdvyrdc8r",
		Amount:  math.NewInt(11600000),
		Claimed: false,
	},
	{
		Address: "elys17e6cqmepnyy3fdfxvw4qkg2h9085x9rhzknz8e",
		Amount:  math.NewInt(11500000),
		Claimed: false,
	},
	{
		Address: "elys1tk3epr965f984zaeaffd98geyj82qcththm704",
		Amount:  math.NewInt(11500000),
		Claimed: false,
	},
	{
		Address: "elys17rqgzc37j9f73hgfjdhzr8ejku6g22w5lex6gs",
		Amount:  math.NewInt(11400000),
		Claimed: false,
	},
	{
		Address: "elys19w9vt8q9cpjay76f784whzmwssvr7ufjp8xwwm",
		Amount:  math.NewInt(11400000),
		Claimed: false,
	},
	{
		Address: "elys1j62xqzl7kjzluu6zr7g2jdx2nayugctf8erere",
		Amount:  math.NewInt(11400000),
		Claimed: false,
	},
	{
		Address: "elys1l6nejehn7ycplvgzsgugemuf7upslm5t0kn2r2",
		Amount:  math.NewInt(11400000),
		Claimed: false,
	},
	{
		Address: "elys1llnudlj42ce5phprtc9v52m0qvx6pu3qx4js87",
		Amount:  math.NewInt(11400000),
		Claimed: false,
	},
	{
		Address: "elys1rfhp67dgh0gmekkq559fs650agj8396vy5gguj",
		Amount:  math.NewInt(11400000),
		Claimed: false,
	},
	{
		Address: "elys1t83nyfq3hsjqlwj6vq3jx5y59el7uxfk5s8gsf",
		Amount:  math.NewInt(11400000),
		Claimed: false,
	},
	{
		Address: "elys1wemjeq22f0zc5zwnfg2wextv6ed7peg29u8my4",
		Amount:  math.NewInt(11400000),
		Claimed: false,
	},
	{
		Address: "elys14m2yvtnpgjv8xskx5uj6alwrukly9y6wwldjyx",
		Amount:  math.NewInt(11300000),
		Claimed: false,
	},
	{
		Address: "elys1vj7uymefmx5s2h3za8e6xs66004k5h36wrmjqk",
		Amount:  math.NewInt(11300000),
		Claimed: false,
	},
	{
		Address: "elys1vrgqdr95mra6ud99spqk6a4k8jzu43n0svlufg",
		Amount:  math.NewInt(11300000),
		Claimed: false,
	},
	{
		Address: "elys105v38cuhcj783d284qu2w8k6m0ddzg25m640dr",
		Amount:  math.NewInt(11200000),
		Claimed: false,
	},
	{
		Address: "elys15qp9mghesvgz6vp4me4j8rude24u4tgtmyy29c",
		Amount:  math.NewInt(11200000),
		Claimed: false,
	},
	{
		Address: "elys19rgz75zd7zclkdenegwex4j92g4a8ntcrczhfp",
		Amount:  math.NewInt(11200000),
		Claimed: false,
	},
	{
		Address: "elys1lugn5z6vthtjx6p46xatkj5770e0m66dwzl6nu",
		Amount:  math.NewInt(11200000),
		Claimed: false,
	},
	{
		Address: "elys120l226z9d3m4p655pn49r8hptwghhf5eezm8v5",
		Amount:  math.NewInt(11000000),
		Claimed: false,
	},
	{
		Address: "elys16rastwuayz800ssjpzdsqh750mw8723advc52j",
		Amount:  math.NewInt(11000000),
		Claimed: false,
	},
	{
		Address: "elys1vuxn7uft9804emt4fp9z2hfeuq4faz4csl9yt0",
		Amount:  math.NewInt(11000000),
		Claimed: false,
	},
	{
		Address: "elys1x9hh2qty2xhmaw3zwzuq4zssy76d67umgqywhu",
		Amount:  math.NewInt(11000000),
		Claimed: false,
	},
	{
		Address: "elys1rp3s43wh9tdghrqn8w7fdesqe0p5a47a7njxxm",
		Amount:  math.NewInt(10900000),
		Claimed: false,
	},
	{
		Address: "elys1q3sp4yxk3epf92sgs2l9tz84543kedpunah69c",
		Amount:  math.NewInt(10800000),
		Claimed: false,
	},
	{
		Address: "elys1yca7ax23q92cjy2wvat05wfq6n5vv27az655q6",
		Amount:  math.NewInt(10800000),
		Claimed: false,
	},
	{
		Address: "elys1qrhch0jsrpapurf99my7hy2a9330ntjpum5fsq",
		Amount:  math.NewInt(10700000),
		Claimed: false,
	},
	{
		Address: "elys1y5ny4xt3afkdmcwa42j3c36f39547th85ghnt3",
		Amount:  math.NewInt(10700000),
		Claimed: false,
	},
	{
		Address: "elys1mnxsyppp0s5x90wu3yfal38h4d7rwqge9qmmg5",
		Amount:  math.NewInt(10600000),
		Claimed: false,
	},
	{
		Address: "elys159249u7q0qx5yj93scqe946g3tkjy722sdm767",
		Amount:  math.NewInt(10500000),
		Claimed: false,
	},
	{
		Address: "elys18rfw8avnkf36l0v852ctshs2f4np2ydlaavaxq",
		Amount:  math.NewInt(10500000),
		Claimed: false,
	},
	{
		Address: "elys1zkll097pmawlsrkdtn9f8x820egzljq3ctx8av",
		Amount:  math.NewInt(10500000),
		Claimed: false,
	},
	{
		Address: "elys1088cyhamer54e275eugwl4na59hcct267f7sk4",
		Amount:  math.NewInt(10400000),
		Claimed: false,
	},
	{
		Address: "elys13fldejsmh0rpyf7fep870r0cx6tjgpfm2a4qqc",
		Amount:  math.NewInt(10400000),
		Claimed: false,
	},
	{
		Address: "elys14gsfqp64efgr0at40ak8ysldr7jpajetfv2cmf",
		Amount:  math.NewInt(10400000),
		Claimed: false,
	},
	{
		Address: "elys1unnl5lxf49v25w4w2eunxv572p486kmlkreamz",
		Amount:  math.NewInt(10400000),
		Claimed: false,
	},
	{
		Address: "elys1yzwgn4c74z3jxsmrgx7nscenjmplehvmdu6dfp",
		Amount:  math.NewInt(10400000),
		Claimed: false,
	},
	{
		Address: "elys1aqct0x3wr2k6gjnflaqq2fr02awk3mc7r6e8mt",
		Amount:  math.NewInt(10300000),
		Claimed: false,
	},
	{
		Address: "elys1hyfjs7s589pwkd4jzs22wmnjuqyxk9cegw4j76",
		Amount:  math.NewInt(10300000),
		Claimed: false,
	},
	{
		Address: "elys1vga7p82pnefdrwtg7736q3dhmx5xdm9mh7eks4",
		Amount:  math.NewInt(10300000),
		Claimed: false,
	},
	{
		Address: "elys1u2f6le28a6hxnchfyhjxr3u7lf5uusgau8jqny",
		Amount:  math.NewInt(10100000),
		Claimed: false,
	},
	{
		Address: "elys12zekp8jylw8s7nvss3zcpa0a92tdghhhzfny4a",
		Amount:  math.NewInt(10000000),
		Claimed: false,
	},
	{
		Address: "elys15kp64e87k8eyc8mxjza4gje0g296u72s2ek8cu",
		Amount:  math.NewInt(10000000),
		Claimed: false,
	},
	{
		Address: "elys1xcd9xxcwllt09uvuf49rq7vdjsn0dzw7ypxaf9",
		Amount:  math.NewInt(10000000),
		Claimed: false,
	},
	{
		Address: "elys1xkz8xwpzkjwru77cvkjjqsgcejcvnrzx3uk0ya",
		Amount:  math.NewInt(10000000),
		Claimed: false,
	},
	{
		Address: "elys1zkguyjnkny2zaj066sed36428h4q9yajvsdns8",
		Amount:  math.NewInt(10000000),
		Claimed: false,
	},
	{
		Address: "elys1nxuzveq885kezvzy4hzsrw4dl5t53d02nv46x8",
		Amount:  math.NewInt(9900000),
		Claimed: false,
	},
	{
		Address: "elys1sqt5d82maw8vlagjjtym84u0ghqe3yvurdc6s0",
		Amount:  math.NewInt(9900000),
		Claimed: false,
	},
	{
		Address: "elys1xst59p2krdwfd2x3falgakrd0yxma0jxtq2yts",
		Amount:  math.NewInt(9900000),
		Claimed: false,
	},
	{
		Address: "elys1gqcn0v3c69zem0j36ds6azq02gdcsynjsdsvu2",
		Amount:  math.NewInt(9800000),
		Claimed: false,
	},
	{
		Address: "elys1mdtrakdkj20hr46r5mmna3f7g59kcdscxqn6td",
		Amount:  math.NewInt(9800000),
		Claimed: false,
	},
	{
		Address: "elys1sskcplfdxv2803wf24rh50kn0vfphs8a7nmses",
		Amount:  math.NewInt(9800000),
		Claimed: false,
	},
	{
		Address: "elys1twy555lpgf9t0d9xf9d5xmd9ee5y2dshjds0ky",
		Amount:  math.NewInt(9800000),
		Claimed: false,
	},
	{
		Address: "elys1aynqjazgtycv60t0y3nplsjmh636wqzg2fvpyt",
		Amount:  math.NewInt(9700000),
		Claimed: false,
	},
	{
		Address: "elys1qxacjj0rdeafflkamm7jws88ft8e2f8ddg9r5m",
		Amount:  math.NewInt(9700000),
		Claimed: false,
	},
	{
		Address: "elys126gjgyg2u22edxamsul80zf93qa834epluwn8a",
		Amount:  math.NewInt(9600000),
		Claimed: false,
	},
	{
		Address: "elys152z6twrlf0dq746kpysp40jp4ws2kue976929g",
		Amount:  math.NewInt(9600000),
		Claimed: false,
	},
	{
		Address: "elys10282v572qt403g677nrulhmep8qyhkcjxug9xc",
		Amount:  math.NewInt(9500000),
		Claimed: false,
	},
	{
		Address: "elys1h5m8skda9n7vwnxma2xmvgnlhcnt3auhx6567r",
		Amount:  math.NewInt(9500000),
		Claimed: false,
	},
	{
		Address: "elys1nrnvcl6h8utpxn4gcy4vzss3k40vy0l9wc5p8r",
		Amount:  math.NewInt(9500000),
		Claimed: false,
	},
	{
		Address: "elys13d6x9hc7y74t8f6yq0mg9mwxsjzzn5vc9lm7ek",
		Amount:  math.NewInt(9300000),
		Claimed: false,
	},
	{
		Address: "elys19f7tw97mpwhzsgezyr8ldaqccg47ylm0vsvnck",
		Amount:  math.NewInt(9300000),
		Claimed: false,
	},
	{
		Address: "elys1a05ys4yd3v9utmc3ks39vtvhgxyyklkjka59ma",
		Amount:  math.NewInt(9200000),
		Claimed: false,
	},
	{
		Address: "elys1gsx8g0uk4nn6lgqzw55fd2mv54dvzskz40zd2g",
		Amount:  math.NewInt(9200000),
		Claimed: false,
	},
	{
		Address: "elys1kgqqv687fqtdjdw7uqvcl5njj9xzuxl6fc506p",
		Amount:  math.NewInt(9200000),
		Claimed: false,
	},
	{
		Address: "elys1kkjq3e5xgjm9ehdvg8kkk5zpcgj8zv7m3d33yh",
		Amount:  math.NewInt(9200000),
		Claimed: false,
	},
	{
		Address: "elys18fc4rdn03x2cq2mdnngmwytqvnlhg9pnem8gv7",
		Amount:  math.NewInt(9100000),
		Claimed: false,
	},
	{
		Address: "elys1hzv5cnxllysqgszr954mg8hhhlsryrcaqn4fry",
		Amount:  math.NewInt(9100000),
		Claimed: false,
	},
	{
		Address: "elys1zek0z7kuavur02u4j0t27d86x2ykshjcmyxhnm",
		Amount:  math.NewInt(9100000),
		Claimed: false,
	},
	{
		Address: "elys13s0jffw3z0rgjp2envcwh3au2hm3qwghdwxjtw",
		Amount:  math.NewInt(9000000),
		Claimed: false,
	},
	{
		Address: "elys19ca0u4dsdfsgq6vqfhaqc08hpjm6vn5fmmuqg7",
		Amount:  math.NewInt(9000000),
		Claimed: false,
	},
	{
		Address: "elys16ywwpukd8rhmkqngdhpwvhu3vhm9lq0qgkr7nv",
		Amount:  math.NewInt(8900000),
		Claimed: false,
	},
	{
		Address: "elys1qh5ucy4f55sttn9h200c22weum7juh0u4vm3fy",
		Amount:  math.NewInt(8900000),
		Claimed: false,
	},
	{
		Address: "elys157zrt6vg3jz88yhlnlzrgfre5w2shdscnt96hg",
		Amount:  math.NewInt(8800000),
		Claimed: false,
	},
	{
		Address: "elys176ewe3ks7thvrzupssfuyu6tcy9yww5pgp40xl",
		Amount:  math.NewInt(8800000),
		Claimed: false,
	},
	{
		Address: "elys1ct8xcqch05jv7rxr5q5k7kt5fr2q0ynygy497m",
		Amount:  math.NewInt(8800000),
		Claimed: false,
	},
	{
		Address: "elys1gygl7d0s2e23xf23a28z6kxja3crmkdrr7gpmk",
		Amount:  math.NewInt(8800000),
		Claimed: false,
	},
	{
		Address: "elys1hueqkp6h9kheahdmqtd4a8quk9qnzqg6gkkug7",
		Amount:  math.NewInt(8800000),
		Claimed: false,
	},
	{
		Address: "elys1nfjfvyk9a4s4rerusgh5v8nmyp0u4yqwamtvs7",
		Amount:  math.NewInt(8800000),
		Claimed: false,
	},
	{
		Address: "elys1u9gahjw3zfkahfj0t4yvsj49paxtvly4u2n5s4",
		Amount:  math.NewInt(8800000),
		Claimed: false,
	},
	{
		Address: "elys1vxud4v2f8l6c89wczt7wrhlrhp2rm2jx3h3y4f",
		Amount:  math.NewInt(8800000),
		Claimed: false,
	},
	{
		Address: "elys1wh7g24hary7g5prqj33fu2wgz76757ytxx98un",
		Amount:  math.NewInt(8800000),
		Claimed: false,
	},
	{
		Address: "elys1wrzr6tsyyf5ljcaqxrqfrs6za25wm6va9j944c",
		Amount:  math.NewInt(8800000),
		Claimed: false,
	},
	{
		Address: "elys1xlug3l0lpusccl6asedx6hspjsrl5aq4rxnc3q",
		Amount:  math.NewInt(8800000),
		Claimed: false,
	},
	{
		Address: "elys1ee55g6450hlw6rjwal6mkmm5q9ftxvct2x7kre",
		Amount:  math.NewInt(8700000),
		Claimed: false,
	},
	{
		Address: "elys1g4qatlj945qm3wufftm8chy60trmjwc43k6ty9",
		Amount:  math.NewInt(8700000),
		Claimed: false,
	},
	{
		Address: "elys12yxlpwjv0apser40dhn49hrtxr3ekt42txtcdn",
		Amount:  math.NewInt(8600000),
		Claimed: false,
	},
	{
		Address: "elys13p4q85rtd7nzx876lax2w7tqds2q0jmygfcfpf",
		Amount:  math.NewInt(8600000),
		Claimed: false,
	},
	{
		Address: "elys14wd4f3sa3r65n3rj25l66fjvepgzyd50pzhn3g",
		Amount:  math.NewInt(8600000),
		Claimed: false,
	},
	{
		Address: "elys1dn6dwu33n0m5gjpw3y746apxkczm6k4pyvyq57",
		Amount:  math.NewInt(8600000),
		Claimed: false,
	},
	{
		Address: "elys17ltxg829zywshwc4lm2l98dmfx4m7my2tzylq4",
		Amount:  math.NewInt(8500000),
		Claimed: false,
	},
	{
		Address: "elys1c5s3zjycnlah708vtjgfnts29m9y9gae2lkhu0",
		Amount:  math.NewInt(8400000),
		Claimed: false,
	},
	{
		Address: "elys1gtl8lmtj30rkyculrgs0hveys8dflly72trfhk",
		Amount:  math.NewInt(8400000),
		Claimed: false,
	},
	{
		Address: "elys1xpvy7wx4v45m3ccapytwe6kjhsv4zjfp3ew2lt",
		Amount:  math.NewInt(8400000),
		Claimed: false,
	},
	{
		Address: "elys1qfpf0mdx93jqaqw5pxurswng22azy4h3lc3g6j",
		Amount:  math.NewInt(8300000),
		Claimed: false,
	},
	{
		Address: "elys1536ywmjll5asqu5ewhquhwscxefxwm0dwz85p5",
		Amount:  math.NewInt(8199999),
		Claimed: false,
	},
	{
		Address: "elys19uvm8x2lqe83qjm5a7j7uvycaec3yyf7ukmc80",
		Amount:  math.NewInt(8199999),
		Claimed: false,
	},
	{
		Address: "elys1nmhm5kufv08en7g3fa2emuzu0xq5gwam4u3hz2",
		Amount:  math.NewInt(8199999),
		Claimed: false,
	},
	{
		Address: "elys1sfdvyjjsp2jv5v2weayhtd4xma3s85htakr2t4",
		Amount:  math.NewInt(8100000),
		Claimed: false,
	},
	{
		Address: "elys14swvl8umzqvya2e83hwy4vs4374v3ndc7z0euv",
		Amount:  math.NewInt(8000000),
		Claimed: false,
	},
	{
		Address: "elys12hf403q9m5dxwpqu5uaxdfw4qqyq3djms89646",
		Amount:  math.NewInt(7900000),
		Claimed: false,
	},
	{
		Address: "elys16vvmjpsduqed58rlju2kl9tes5zpznexahxzk5",
		Amount:  math.NewInt(7900000),
		Claimed: false,
	},
	{
		Address: "elys17r99kp7y6fat7d2466vyq7fthq9y26cyclwq3w",
		Amount:  math.NewInt(7900000),
		Claimed: false,
	},
	{
		Address: "elys1trut9nguyyeznten8lhv349zqa0aazz8sappf3",
		Amount:  math.NewInt(7900000),
		Claimed: false,
	},
	{
		Address: "elys1zts4j8ulets0mwaw5kky7f9aly7ce4thf9npm9",
		Amount:  math.NewInt(7900000),
		Claimed: false,
	},
	{
		Address: "elys10jrz387jvyyqggfvyu07j800cvtljr032exvy0",
		Amount:  math.NewInt(7800000),
		Claimed: false,
	},
	{
		Address: "elys1xdumrvfr8cmv667r7tqytnf5dlrx84vg596qa0",
		Amount:  math.NewInt(7700000),
		Claimed: false,
	},
	{
		Address: "elys166srl2fqww5kutmwkf4swapmztc8y4wnyw3lct",
		Amount:  math.NewInt(7600000),
		Claimed: false,
	},
	{
		Address: "elys18rzv288h3pj9v66egjm0wn3aezz9gkhpp3q0hg",
		Amount:  math.NewInt(7600000),
		Claimed: false,
	},
	{
		Address: "elys1d6h3qgvyn5ef5x4ke35hpxwhgwq0za5azv9sm7",
		Amount:  math.NewInt(7600000),
		Claimed: false,
	},
	{
		Address: "elys1y98l7mqjzyeuc9was4d8p5lsrj8vnsdakatrw2",
		Amount:  math.NewInt(7600000),
		Claimed: false,
	},
	{
		Address: "elys14za6x5yzftk4n5mjtr42acfjef77v5akqal70j",
		Amount:  math.NewInt(7500000),
		Claimed: false,
	},
	{
		Address: "elys1t74ks4ex55zx04wlnx9hauphgzn0qsc2rkazh6",
		Amount:  math.NewInt(7500000),
		Claimed: false,
	},
	{
		Address: "elys16mhftrhncgvp9j6u3zz87qc69n3e3ha6mud3sg",
		Amount:  math.NewInt(7400000),
		Claimed: false,
	},
	{
		Address: "elys17pnw7af62uhftue9lc87q3l4s3y2mply5kj5mx",
		Amount:  math.NewInt(7400000),
		Claimed: false,
	},
	{
		Address: "elys1cxtltr655vy0fznd8ctexkxk6cppsuyjdzhske",
		Amount:  math.NewInt(7400000),
		Claimed: false,
	},
	{
		Address: "elys1gnxdak8edldkhsmkl93ez0d39qpyrkzahq2att",
		Amount:  math.NewInt(7400000),
		Claimed: false,
	},
	{
		Address: "elys1va4hpxzy4ffn9cuwtg3lw3cagqr55crvrkwazm",
		Amount:  math.NewInt(7400000),
		Claimed: false,
	},
	{
		Address: "elys1xp2rudyuz6sw53j44nltrrr42r9uzv4kdr6dy0",
		Amount:  math.NewInt(7400000),
		Claimed: false,
	},
	{
		Address: "elys1el583l7u5trln9glg8uv2uk5eys7g7h3xfu0df",
		Amount:  math.NewInt(7300000),
		Claimed: false,
	},
	{
		Address: "elys1p404vzfjhm9z0vxz68uj6ypcqrn78lfe0avkwl",
		Amount:  math.NewInt(7300000),
		Claimed: false,
	},
	{
		Address: "elys1qypvnc2mxnn5jxdd7e8wa6ydldeqg2vynrz2ll",
		Amount:  math.NewInt(7300000),
		Claimed: false,
	},
	{
		Address: "elys1wkeu28y0al4k8vpa9rg7qpw5qf2yx9x4nsc6h4",
		Amount:  math.NewInt(7300000),
		Claimed: false,
	},
	{
		Address: "elys1zf298fth8yu7czq3eyyh4g90n3sc5e2ksjqwva",
		Amount:  math.NewInt(7300000),
		Claimed: false,
	},
	{
		Address: "elys10f482htpdaahnsh2kch3sl8kjc322rc43y5fmp",
		Amount:  math.NewInt(7200000),
		Claimed: false,
	},
	{
		Address: "elys1avmvkag66ltqgl7r0q7s7te3jds4lacg82un9c",
		Amount:  math.NewInt(7200000),
		Claimed: false,
	},
	{
		Address: "elys1nmxzeaccjqh4t43lgkmh63ujfunjaq7dnwmsv5",
		Amount:  math.NewInt(7200000),
		Claimed: false,
	},
	{
		Address: "elys1y2ducl7a7mal56vjh9a0g3d3dcnx93k8850ck9",
		Amount:  math.NewInt(7200000),
		Claimed: false,
	},
	{
		Address: "elys1jtyrduppjjzqkw6twtk8fnrts3chzxu3rhd0kl",
		Amount:  math.NewInt(7100000),
		Claimed: false,
	},
	{
		Address: "elys1s5zvxsz0gvjes4tmugdv92jphqe0g0lpnws5k9",
		Amount:  math.NewInt(7100000),
		Claimed: false,
	},
	{
		Address: "elys1uhqw85rtvj4r4y7u6cnz4rhjaewq0fmw88ghdk",
		Amount:  math.NewInt(7100000),
		Claimed: false,
	},
	{
		Address: "elys1vhntxhc735mys7zftxmlkquts9ca2x8gzsx948",
		Amount:  math.NewInt(7100000),
		Claimed: false,
	},
	{
		Address: "elys1ysrgc7rwpdkytlwwqj4zqtz26cz8pvvdsfgx3s",
		Amount:  math.NewInt(7100000),
		Claimed: false,
	},
	{
		Address: "elys1z660r3m82ufahzfwsp7r59h2rjv5ryad0p5chc",
		Amount:  math.NewInt(7100000),
		Claimed: false,
	},
	{
		Address: "elys1k4cu9rrkpe65fe3p9xq53nn8w2tm5mns52n5ur",
		Amount:  math.NewInt(7000000),
		Claimed: false,
	},
	{
		Address: "elys1l5mmjndyneee6s862pl4s00l989dwyp4vrngsj",
		Amount:  math.NewInt(7000000),
		Claimed: false,
	},
	{
		Address: "elys1efhqf7c8k630uvz5pa52np8zfydda3m7n4ru4w",
		Amount:  math.NewInt(6900000),
		Claimed: false,
	},
	{
		Address: "elys1rt45e5mpzyffngafwkh6n2qtccz9x82dgh55f7",
		Amount:  math.NewInt(6900000),
		Claimed: false,
	},
	{
		Address: "elys19lna60jtcyv689l6r4zcp3efxcvljrvelamkzr",
		Amount:  math.NewInt(6800000),
		Claimed: false,
	},
	{
		Address: "elys1djnt97ggyqg2eyn9h8ht66ev03jj3medqjfx2w",
		Amount:  math.NewInt(6800000),
		Claimed: false,
	},
	{
		Address: "elys1ghyx9q63c6m3cew363fwqpjuxzc7c7pd602wcr",
		Amount:  math.NewInt(6800000),
		Claimed: false,
	},
	{
		Address: "elys1nt6et0h7n79h4rh7h9w82hlpjt34usf64mdkvk",
		Amount:  math.NewInt(6800000),
		Claimed: false,
	},
	{
		Address: "elys1yhwdlajgyx2f4d6633yzq8ett00t347u2t9vxr",
		Amount:  math.NewInt(6800000),
		Claimed: false,
	},
	{
		Address: "elys12e7x392uux4xx3suzf7xmhnzdmwalzum49pvc0",
		Amount:  math.NewInt(6700000),
		Claimed: false,
	},
	{
		Address: "elys13nhjpg5x05k6xckphrxhu7982p4sscscltuyc3",
		Amount:  math.NewInt(6700000),
		Claimed: false,
	},
	{
		Address: "elys14eyxfxrvevty00j53ly4h7em53gfrsvl90sejy",
		Amount:  math.NewInt(6700000),
		Claimed: false,
	},
	{
		Address: "elys14lpq3dp7vgtf09dx4j3e8x6heyt66yu73q6yzz",
		Amount:  math.NewInt(6700000),
		Claimed: false,
	},
	{
		Address: "elys19whz3eg6ugkz5uk8e0cuz5yumr57dsj37dq03u",
		Amount:  math.NewInt(6700000),
		Claimed: false,
	},
	{
		Address: "elys1fz707jw2zkew5h7k3rjc69zaektmk85tups8g4",
		Amount:  math.NewInt(6700000),
		Claimed: false,
	},
	{
		Address: "elys1jt6p3s85w4d3x4tent35ep80grch7nwg6yeh2n",
		Amount:  math.NewInt(6700000),
		Claimed: false,
	},
	{
		Address: "elys1pn79gqy0h69adn4myn4npqfvcppsyykkt2adha",
		Amount:  math.NewInt(6700000),
		Claimed: false,
	},
	{
		Address: "elys1r3y4zrwpg8l76mv0cffrdqsw8uv5nstxymym7x",
		Amount:  math.NewInt(6700000),
		Claimed: false,
	},
	{
		Address: "elys1wrxa99mftkuykglup25sc5fdsycx6squqe4379",
		Amount:  math.NewInt(6700000),
		Claimed: false,
	},
	{
		Address: "elys1y0tez5sj0yvscgcktjs923rhmggyxnt88g0czn",
		Amount:  math.NewInt(6700000),
		Claimed: false,
	},
	{
		Address: "elys1udsr9lwnw89n265gjnjg25nmf3qtvlmn3vk4yd",
		Amount:  math.NewInt(6600000),
		Claimed: false,
	},
	{
		Address: "elys14r5960te5hfvrkn3207stjjqymhz74ua8j89qx",
		Amount:  math.NewInt(6500000),
		Claimed: false,
	},
	{
		Address: "elys15j9knpst3rxg8szul6g292p99dtfvlma5dhavg",
		Amount:  math.NewInt(6500000),
		Claimed: false,
	},
	{
		Address: "elys1er42f5rh26nskeze8y2gm32v7t2fu8x8l2kdsj",
		Amount:  math.NewInt(6500000),
		Claimed: false,
	},
	{
		Address: "elys1pq0p9mz0kkrmxdrkqdymlwpeefqdz3zs40hzay",
		Amount:  math.NewInt(6500000),
		Claimed: false,
	},
	{
		Address: "elys1q0k9zrgvs8hzmrv6gsm80r4n4s9xhdjxk9gcq0",
		Amount:  math.NewInt(6500000),
		Claimed: false,
	},
	{
		Address: "elys1wek0l09tgtc9qta0rqjqm4q9qlnxkuux5ywdve",
		Amount:  math.NewInt(6500000),
		Claimed: false,
	},
	{
		Address: "elys12de599hdl8wy3kyqcs7xxhljqcemlv8cm5k09x",
		Amount:  math.NewInt(6400000),
		Claimed: false,
	},
	{
		Address: "elys13fjpd6zexn4v6czuunqzry3ecu3pqel8ude8k2",
		Amount:  math.NewInt(6400000),
		Claimed: false,
	},
	{
		Address: "elys16m4yxukw8sgd0k7w3mwa29kqejq4ym6xzh4cll",
		Amount:  math.NewInt(6400000),
		Claimed: false,
	},
	{
		Address: "elys1qctsnew9jprsk2jdngu6ddemgd3h4f49cnlj6h",
		Amount:  math.NewInt(6400000),
		Claimed: false,
	},
	{
		Address: "elys1hd5jg34ez2wepvwrj97d8lqu45hxqjjqkhwcma",
		Amount:  math.NewInt(6300000),
		Claimed: false,
	},
	{
		Address: "elys1kqzfzpzja4jerjxpqadyq0vdf750ku9uvktwrj",
		Amount:  math.NewInt(6300000),
		Claimed: false,
	},
	{
		Address: "elys1qg9c05ce5puhfz838u7ljj57nfgrn3v6u6l2ex",
		Amount:  math.NewInt(6300000),
		Claimed: false,
	},
	{
		Address: "elys1t2xl92qmc6lppkn6ys4qryjzrj5jwhjwqsewjt",
		Amount:  math.NewInt(6300000),
		Claimed: false,
	},
	{
		Address: "elys10vd82mjdlxjk28ue6k4frn8vaj0kuss47gc6as",
		Amount:  math.NewInt(6200000),
		Claimed: false,
	},
	{
		Address: "elys1406rw46yzzl32azkrnzc2gagc48ypwerpk6uck",
		Amount:  math.NewInt(6200000),
		Claimed: false,
	},
	{
		Address: "elys1clrsfkyjcmfzxfukdcnktk04q7l7aht8p4epug",
		Amount:  math.NewInt(6200000),
		Claimed: false,
	},
	{
		Address: "elys1l6zt73u4499gry8pscynxktlljchexzx5x9c3s",
		Amount:  math.NewInt(6200000),
		Claimed: false,
	},
	{
		Address: "elys1lvw59tqg40w5qdxvardpkqs0yexxhvdpw9cl22",
		Amount:  math.NewInt(6200000),
		Claimed: false,
	},
	{
		Address: "elys130v75mvnf7yarhnxhnd6pktjdu86psyregazry",
		Amount:  math.NewInt(6100000),
		Claimed: false,
	},
	{
		Address: "elys1jwpxyu529y0edsa0p56qnp7ez3f260sqp4h2r8",
		Amount:  math.NewInt(6100000),
		Claimed: false,
	},
	{
		Address: "elys1r44dspqmc25fe0rp3z8tztfmzhexrekg38w8js",
		Amount:  math.NewInt(6100000),
		Claimed: false,
	},
	{
		Address: "elys167uwkfmz355dp8t02jrgs4gkrh7v605dg77q20",
		Amount:  math.NewInt(6000000),
		Claimed: false,
	},
	{
		Address: "elys176lj0acmx535p6n6ymku40u2ekjcn2epenlfum",
		Amount:  math.NewInt(6000000),
		Claimed: false,
	},
	{
		Address: "elys1fp0zhnamluxapeyqhnnh03qhurf5v3sfs9yu2f",
		Amount:  math.NewInt(6000000),
		Claimed: false,
	},
	{
		Address: "elys132r0uejcgpgu3y5ys2r7888tw90g7twyjaen4q",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys13cwec303gkvykwpfu03x4ur53vgssld3a5mf30",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys18rkez04fxgqnwdv8s4dd4fteyfm0ym6esv3llv",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1amys7y4kgtgxcsgnnz0umg4f7056xg48mn3kxj",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1anvcap7tvcnna7vxlvgjquu0g8lupkthep04x6",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1k837f230hhu9ycvh6rv89ga7cvjfnwlfwcsge5",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1mpf3hpe2j5mvqdv89vm7wus02feqhu3pzrzstn",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1s8jtdrajecjt0d7g0d3zypd2gnaq7wjg7rak3f",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1trl4acgmvlh8zwf8nj6uwt0vv97wevaqklewmp",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1u2dce72xvrd048ph9f49e5jjy9uv2jhtvqpqul",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1v6qgz5qjv67pm59pwp2mxl384va7zpm6axn8as",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1z7ghy2l39juz6cperrgrr7rme5j9mnd7vjsjrh",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1zsy895gz2pe0ux0alx07a88l485qlg6tg35nep",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1zyqy5l4xpw3s9ytw4nu9cymc8rt4c9kjknhh64",
		Amount:  math.NewInt(5900000),
		Claimed: false,
	},
	{
		Address: "elys1m0d9a3t5kn9myecpvu7kleggv8tg0mruntr4uy",
		Amount:  math.NewInt(5800000),
		Claimed: false,
	},
	{
		Address: "elys1x867jl0dktma6fasl9jj399vcrgrkuz2tp3jnp",
		Amount:  math.NewInt(5800000),
		Claimed: false,
	},
	{
		Address: "elys1yr58hk43p8fhq3fdsaa4kw37pg39wqlugluq46",
		Amount:  math.NewInt(5800000),
		Claimed: false,
	},
	{
		Address: "elys125ge8tpryflx28xf05kh3vwdllthxtesya6gq4",
		Amount:  math.NewInt(5700000),
		Claimed: false,
	},
	{
		Address: "elys127npkqhnwhywsmqwmpjzku2ae0wf4em9uptpmd",
		Amount:  math.NewInt(5700000),
		Claimed: false,
	},
	{
		Address: "elys16tu35fvap46znlpk2jag7y2j9uyysjr6nfxrc9",
		Amount:  math.NewInt(5700000),
		Claimed: false,
	},
	{
		Address: "elys19u2fd64fq2qds3klk5czzjukq4tva09az5cs05",
		Amount:  math.NewInt(5700000),
		Claimed: false,
	},
	{
		Address: "elys1ml4vzagl8mkfaywk822ld49ht4v20lfve9z6ef",
		Amount:  math.NewInt(5700000),
		Claimed: false,
	},
	{
		Address: "elys14j627e8n4pkpttweawqacwvwq433ylkn68kesp",
		Amount:  math.NewInt(5600000),
		Claimed: false,
	},
	{
		Address: "elys177d58ak58xap226u7a8xm02qdvhuh2jghzdufy",
		Amount:  math.NewInt(5600000),
		Claimed: false,
	},
	{
		Address: "elys1jtwjajt76k33krkzswtsymand6d5jdh4p9kr47",
		Amount:  math.NewInt(5600000),
		Claimed: false,
	},
	{
		Address: "elys1ug2hectpjkxg20xreze6nzhkp67enp6kh64mv9",
		Amount:  math.NewInt(5600000),
		Claimed: false,
	},
	{
		Address: "elys1ulxqaxhuld897wmnlt8vdlm5d2vcaf6qakafxw",
		Amount:  math.NewInt(5600000),
		Claimed: false,
	},
	{
		Address: "elys18ned0nh5d9qdnjvv9jtdz26mka3guuwa8yh8ry",
		Amount:  math.NewInt(5500000),
		Claimed: false,
	},
	{
		Address: "elys1elewysfgr48ulpc4e0lrt2e934llajknurpxjr",
		Amount:  math.NewInt(5400000),
		Claimed: false,
	},
	{
		Address: "elys1t099ka67ly7ge3rccmumc08tqkp8685539379q",
		Amount:  math.NewInt(5400000),
		Claimed: false,
	},
	{
		Address: "elys165qdr7j4cjjpl5ghqmzjt4jpyk0fja6tqvkucx",
		Amount:  math.NewInt(5200000),
		Claimed: false,
	},
	{
		Address: "elys16xl4rzle83n048t3z2md3ld9yxzvl0nud3lyzx",
		Amount:  math.NewInt(5200000),
		Claimed: false,
	},
	{
		Address: "elys14rh6hp8fy3hls692yns5lr56z9976rc3smxzcl",
		Amount:  math.NewInt(5100000),
		Claimed: false,
	},
	{
		Address: "elys1ajueefk0p3jxygxf0qqtkww537ltjgc4lpzanh",
		Amount:  math.NewInt(5100000),
		Claimed: false,
	},
	{
		Address: "elys1frrlnxh5advkx934as6wwlwf0ekxntc7p4jvfn",
		Amount:  math.NewInt(5100000),
		Claimed: false,
	},
	{
		Address: "elys1w96erxzuvvss4s296v8amn9jadeyx442j3j7ea",
		Amount:  math.NewInt(5100000),
		Claimed: false,
	},
	{
		Address: "elys15cr6kv94v5q7ca5290nrwv9c7suqex6t66a7tz",
		Amount:  math.NewInt(5000000),
		Claimed: false,
	},
	{
		Address: "elys15d67x8akwf4amnuhlcx63tv4s273rlur73zxsq",
		Amount:  math.NewInt(5000000),
		Claimed: false,
	},
	{
		Address: "elys1ss4krazunlfnc8munks4uvsa5c42vlwr2smejj",
		Amount:  math.NewInt(5000000),
		Claimed: false,
	},
	{
		Address: "elys1vzcu8wyrmrkjvreje5c8w6tn7a4x3g6yzrmwyy",
		Amount:  math.NewInt(5000000),
		Claimed: false,
	},
	{
		Address: "elys12x6dkswewgzm673d9q2qsdhvkdra92dcp46h37",
		Amount:  math.NewInt(4900000),
		Claimed: false,
	},
	{
		Address: "elys14gatzexre0d29dwmhvzmr37qvqn0c5dzt9zn5j",
		Amount:  math.NewInt(4900000),
		Claimed: false,
	},
	{
		Address: "elys1jxryl4k5zt9emqhe2kgy50pheaeut94er69aly",
		Amount:  math.NewInt(4900000),
		Claimed: false,
	},
	{
		Address: "elys1p59rhlurlhm2qrnyeflg2dwnzdv9cxqd7sz0cq",
		Amount:  math.NewInt(4900000),
		Claimed: false,
	},
	{
		Address: "elys1p7k8f0erllcq8hsmxzxmt4we97gwp0aa959rrw",
		Amount:  math.NewInt(4900000),
		Claimed: false,
	},
	{
		Address: "elys1pa0q7vmnm8x26rpua4lrhnga8uww0hrcu8xwwt",
		Amount:  math.NewInt(4900000),
		Claimed: false,
	},
	{
		Address: "elys1rku6z5u3vtdr7hcnsdqdj8e79c2ahzxz3r78fa",
		Amount:  math.NewInt(4900000),
		Claimed: false,
	},
	{
		Address: "elys1vfa77858lr8klnkek48v3slfjmjck9pnpsftyx",
		Amount:  math.NewInt(4900000),
		Claimed: false,
	},
	{
		Address: "elys1xl82fj5k392qgcwsp4tvktu9evj6y0204hs4tg",
		Amount:  math.NewInt(4900000),
		Claimed: false,
	},
	{
		Address: "elys1m643etr6xu0j9jsre6mwrhf7sxkn443dvalpn6",
		Amount:  math.NewInt(4800000),
		Claimed: false,
	},
	{
		Address: "elys1mq6ltf08rud3fp0h7glcv45z5d768a9dpsvda6",
		Amount:  math.NewInt(4800000),
		Claimed: false,
	},
	{
		Address: "elys1nju4s032g6ecxel6xu2w6fd7a4vy0nu0vf3wx6",
		Amount:  math.NewInt(4800000),
		Claimed: false,
	},
	{
		Address: "elys1e2t2zqhl2nmsj3wap9k2svlr6m4rh2w82nt07v",
		Amount:  math.NewInt(4700000),
		Claimed: false,
	},
	{
		Address: "elys1fdelwf0glm24kwkl0hla3vzd6qxaay7s4tx9rl",
		Amount:  math.NewInt(4700000),
		Claimed: false,
	},
	{
		Address: "elys10gev5rcvh8p57p26fwrtdenx6qx0ywsdtx6hyy",
		Amount:  math.NewInt(4600000),
		Claimed: false,
	},
	{
		Address: "elys1czmj52v2mx6r3vdny4xy6v8fxsss50j28yrlar",
		Amount:  math.NewInt(4600000),
		Claimed: false,
	},
	{
		Address: "elys1e7gyw75l7vw4hd5wutgq2mgmrtw0ch0msqujv0",
		Amount:  math.NewInt(4600000),
		Claimed: false,
	},
	{
		Address: "elys1j8cn8vmqe0y42alulzf9zaqx8euytml4eapvth",
		Amount:  math.NewInt(4600000),
		Claimed: false,
	},
	{
		Address: "elys1rhptyjt80zlkkp56k7j3lzujh3j9f2v0jt0v8h",
		Amount:  math.NewInt(4600000),
		Claimed: false,
	},
	{
		Address: "elys1ucexv9agj0zhd2dffkg006l3jzaxxgfddl33m5",
		Amount:  math.NewInt(4600000),
		Claimed: false,
	},
	{
		Address: "elys1xrpg6ecelvcaxe90f64wtm3q4mcksw655kepa9",
		Amount:  math.NewInt(4600000),
		Claimed: false,
	},
	{
		Address: "elys1z8u0hxv90u0gtvfjvf9q82l8laflg62nr85ytz",
		Amount:  math.NewInt(4600000),
		Claimed: false,
	},
	{
		Address: "elys14xnk3gr40jqlu07arecrmr73s89w8emtscgxse",
		Amount:  math.NewInt(4500000),
		Claimed: false,
	},
	{
		Address: "elys15qcjj09nv0lp866p5f5wel39qxknglnxml0q68",
		Amount:  math.NewInt(4500000),
		Claimed: false,
	},
	{
		Address: "elys190vcdpegef4m4qmxwpnq6u2wzkyzsthgc7w6ec",
		Amount:  math.NewInt(4500000),
		Claimed: false,
	},
	{
		Address: "elys1985q5wtq8kr2aus5sr5h26pru8cdeayq8flrw9",
		Amount:  math.NewInt(4500000),
		Claimed: false,
	},
	{
		Address: "elys1e8k9thgwch6kpdw2eamqzyex93dqjy07psqw0p",
		Amount:  math.NewInt(4500000),
		Claimed: false,
	},
	{
		Address: "elys1m3ua5hzh656q2dt7q90cv7vrdknavmyancppkp",
		Amount:  math.NewInt(4500000),
		Claimed: false,
	},
	{
		Address: "elys1rvzm7aptrv7ustwyy659n4z6aptz563fx9a8d6",
		Amount:  math.NewInt(4500000),
		Claimed: false,
	},
	{
		Address: "elys1tevjgyj2z5trs4yxqxjdf2yl2y88c2gn0wrehq",
		Amount:  math.NewInt(4500000),
		Claimed: false,
	},
	{
		Address: "elys1uackw5kc9sqtjh79c339xz98juv77ndr3ukscv",
		Amount:  math.NewInt(4500000),
		Claimed: false,
	},
	{
		Address: "elys1uq34ur7juqlfw5s8l0aa6w9cj5pll6yastx9df",
		Amount:  math.NewInt(4500000),
		Claimed: false,
	},
	{
		Address: "elys107umxemsupzedafex36mg6et48rcc0n0337m5j",
		Amount:  math.NewInt(4400000),
		Claimed: false,
	},
	{
		Address: "elys13uxcjgaxyawd2vtxy4kacd4xza9urmlnclgnar",
		Amount:  math.NewInt(4400000),
		Claimed: false,
	},
	{
		Address: "elys15t87aqrfflny67k45mln82qzf88dsantjsz2md",
		Amount:  math.NewInt(4400000),
		Claimed: false,
	},
	{
		Address: "elys1hkrtpeefcjrcd934gc3s03qnhk0ma2eun3q37e",
		Amount:  math.NewInt(4400000),
		Claimed: false,
	},
	{
		Address: "elys1ypk8p72lqfezjpypxwu6n69de6qzm2sylrtya8",
		Amount:  math.NewInt(4400000),
		Claimed: false,
	},
	{
		Address: "elys1c26gemzgtump82gw26wxxrtzr79mpxtv8yng68",
		Amount:  math.NewInt(4300000),
		Claimed: false,
	},
	{
		Address: "elys1ey6fkh5w8rylvdfujxaa66mraz9dv7eg0pkk3j",
		Amount:  math.NewInt(4300000),
		Claimed: false,
	},
	{
		Address: "elys1f6n30pj0c3hzemf04n7xvczqm40hqj247l33z7",
		Amount:  math.NewInt(4300000),
		Claimed: false,
	},
	{
		Address: "elys1su2nd2ha5efzz7lwvgrkptk7lkuqn357dvzrm4",
		Amount:  math.NewInt(4300000),
		Claimed: false,
	},
	{
		Address: "elys1xs0xfumfg8p8z06ges38vu48e375s7rh9cfs5q",
		Amount:  math.NewInt(4300000),
		Claimed: false,
	},
	{
		Address: "elys1208kwwdndd6ru6qanvf0sh4w9jqqgv0arykkan",
		Amount:  math.NewInt(4200000),
		Claimed: false,
	},
	{
		Address: "elys164dpffm974ayj7j9r33t4yhv7a6wu59q6dnlh5",
		Amount:  math.NewInt(4200000),
		Claimed: false,
	},
	{
		Address: "elys1ag4eggqs4uwhxknx5ycqwc9nu9c2x8l4rgp7h2",
		Amount:  math.NewInt(4200000),
		Claimed: false,
	},
	{
		Address: "elys1hwwrk6gv4nnp6e4h3m7zxn0lys90tnt9m6lyur",
		Amount:  math.NewInt(4200000),
		Claimed: false,
	},
	{
		Address: "elys1w0yh5358q8rgp3z0mdqjq5y5lcjen8yqu65re3",
		Amount:  math.NewInt(4200000),
		Claimed: false,
	},
	{
		Address: "elys1xah935sma5w8ytzmdsj75gcyrnwqct8r9jtfvq",
		Amount:  math.NewInt(4200000),
		Claimed: false,
	},
	{
		Address: "elys1xsy9pxr0g7txf3k4qj3zmdj7ms46av8t8rxpaz",
		Amount:  math.NewInt(4200000),
		Claimed: false,
	},
	{
		Address: "elys1yv089j8w8e9nxlmvhlrg2ak7h382tcxxxt4ttp",
		Amount:  math.NewInt(4200000),
		Claimed: false,
	},
	{
		Address: "elys14rcvw3jma0x8w8qvck24zsvw7ufctf4cq2qwar",
		Amount:  math.NewInt(4099999),
		Claimed: false,
	},
	{
		Address: "elys1h0t3m57m6q88tasvdm6lh0h9tnd9aysjm4y288",
		Amount:  math.NewInt(4099999),
		Claimed: false,
	},
	{
		Address: "elys1pk57n3jxrtts4dsagpuaxg4ptl4pe9ux7wht8j",
		Amount:  math.NewInt(4099999),
		Claimed: false,
	},
	{
		Address: "elys1tw45lkckty8k3gx83lgr03lx23ajm6ay8y57pf",
		Amount:  math.NewInt(4099999),
		Claimed: false,
	},
	{
		Address: "elys1w626x2g777n6xk9nxp4tk8hwne3pc7l9kcdsre",
		Amount:  math.NewInt(4099999),
		Claimed: false,
	},
	{
		Address: "elys1yn4fswg4zyxzxnfs7kv9qggczage0wmwvnylga",
		Amount:  math.NewInt(4099999),
		Claimed: false,
	},
	{
		Address: "elys13let8xcwfcmupuy9kxxd8pg95m7flv5hmr0ywu",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys1686xkke5f8agp9q5pfq646xhykyq7jhe5qczms",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys19s58gmd8ksetqqsrexx9f4gxscu7pmc0l4rqnr",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys1a0unrcj9cw4893z66j22rkx9zltxvjs3c0nmxk",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys1a2673py5cpqr5jvhquxvn7zracaaj0pfqkr833",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys1jfudfx6vrxtl70lfejhn5dy2g7qmlwm7kskcre",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys1l4ll05ug2cswddqnathnjhv6a4ssjq7fkj4fmy",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys1qpeqe8h0yvt9e6q2nc2ws3htrdtjzsm77fxcyc",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys1t6h624mmhutrt03sycfdm9w3yt5zezxn735l2p",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys1tcdq05ezy8d3aqhmacvzjhgvhpddufhqkmg7dc",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys1tr6s9p6nrlpf7sqsxjpr4rvh256029dhfsr0aq",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys1unapcqardx9f4k2tdlhm5tyhyassytwv5xf9zw",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys1xjvq8pypj65s6k6esqf8ckvx7xhnqesrad3958",
		Amount:  math.NewInt(4000000),
		Claimed: false,
	},
	{
		Address: "elys15e7jv2p2wh0rl4x9jct0md5xpx6dszv0a9pe6q",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys15wfx8085u6g29pxpm33js29flmnhz2hgytzynq",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys16hud85alk30s76cwn38ymskk3sx075nzenge56",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys16ulyffeumsa5ukg35y2vjh9zafp9lze8xdm04u",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys1av5kje2f55pdzq600pe5p9aq90j7978xyees9d",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys1eq2ncte6chhe0hk2dpvxjp32ax95z063rpfn2z",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys1lwxcn8a3ntqpdtx0vd4wx3harnx08enqad9033",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys1phgv0htl0ztwwrekshvk52lp3e209e5s0rrv3q",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys1u0jszvy62g72qpe5x3jnwsar9xejf772h5llrk",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys1v0400qmdkhvp2es8kfr5trhehvjynw22eappz3",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys1w5utzyjtn6a66g4d57gzrueqtn7lt63pq58k5n",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys1x6ysr0kx3vuuc5q7u7r9zkjnqtqcrt0ac99fcv",
		Amount:  math.NewInt(3900000),
		Claimed: false,
	},
	{
		Address: "elys13ctuhv64eeg0qdq23uyuxk9tjym7mz2t6c9t6u",
		Amount:  math.NewInt(3800000),
		Claimed: false,
	},
	{
		Address: "elys163e86z6pd7k9y2hxl752xfu4yp72uzfv54zzec",
		Amount:  math.NewInt(3800000),
		Claimed: false,
	},
	{
		Address: "elys1c4qxujfs58x6033s0kdc5xrqvnmg4sl5tu3utz",
		Amount:  math.NewInt(3800000),
		Claimed: false,
	},
	{
		Address: "elys1rudl9ygrudq66lwkf3j02thecfkmzkvf6wt5h9",
		Amount:  math.NewInt(3800000),
		Claimed: false,
	},
	{
		Address: "elys1504h94q0p558fwxwj5sqx4aahul03xx3wmpr97",
		Amount:  math.NewInt(3700000),
		Claimed: false,
	},
	{
		Address: "elys16lmdwdj963ep7h0tvtdsmfet5e0tp7z56r0kpk",
		Amount:  math.NewInt(3700000),
		Claimed: false,
	},
	{
		Address: "elys1ywh0exa3n8eks54k8am4hm7ace5phcm3whm4wf",
		Amount:  math.NewInt(3600000),
		Claimed: false,
	},
	{
		Address: "elys19vgv7kqczujh76rvnqhpfg5ara3ekkth78y46v",
		Amount:  math.NewInt(3500000),
		Claimed: false,
	},
	{
		Address: "elys1a2uwq0q0tyf3xdlmxl9caxcp6qcfrwfjgq9g60",
		Amount:  math.NewInt(3500000),
		Claimed: false,
	},
	{
		Address: "elys1acmc8f7g9eqs6lvhy3m9889d9djze4cz82te5p",
		Amount:  math.NewInt(3500000),
		Claimed: false,
	},
	{
		Address: "elys1smy988kjlaf0p5z0c4axmcdlftk2ufkavjq739",
		Amount:  math.NewInt(3500000),
		Claimed: false,
	},
	{
		Address: "elys1u97hgpwpnusanzjvkqmxxcqle2wpekzfrdeclw",
		Amount:  math.NewInt(3500000),
		Claimed: false,
	},
	{
		Address: "elys1vg2umpm6qyzvct74mjn3x0fnv64vtg4q3grfpw",
		Amount:  math.NewInt(3500000),
		Claimed: false,
	},
	{
		Address: "elys13v9c0slnr6r7wreacxyr0d9f33qvqc9stht5gl",
		Amount:  math.NewInt(3400000),
		Claimed: false,
	},
	{
		Address: "elys1jjut2tkl7d6pk7tv4gjsntjlp0q5vaz9hw2fa7",
		Amount:  math.NewInt(3400000),
		Claimed: false,
	},
	{
		Address: "elys1lf7vq45ez4l2s7ufuyx4h7g6cn2jg96lm2q4kp",
		Amount:  math.NewInt(3400000),
		Claimed: false,
	},
	{
		Address: "elys1tcl6mkalv7ckud9z6epnt70jnnf8mddrnjhcgc",
		Amount:  math.NewInt(3400000),
		Claimed: false,
	},
	{
		Address: "elys1u8hntnh6przrd0ew6eptz7slx9lzvwtzvjjkzk",
		Amount:  math.NewInt(3400000),
		Claimed: false,
	},
	{
		Address: "elys1wmk48yuhvjgdgrd6wxk2w7tdfu4mlvj497069m",
		Amount:  math.NewInt(3400000),
		Claimed: false,
	},
	{
		Address: "elys19kydxjmh6pfwsf87llxlkkmz29ft3dsmnzdgaz",
		Amount:  math.NewInt(3300000),
		Claimed: false,
	},
	{
		Address: "elys1atp6skvrfgzf9exqghfft8cvj8edhxjcn854xq",
		Amount:  math.NewInt(3300000),
		Claimed: false,
	},
	{
		Address: "elys1awymlr02twz2x2z7y2p79twtg8l2pa89xk0w9k",
		Amount:  math.NewInt(3300000),
		Claimed: false,
	},
	{
		Address: "elys1havwfgnrk2ttl4fn2ej9qtvs2n4ha7ufaqz8kf",
		Amount:  math.NewInt(3300000),
		Claimed: false,
	},
	{
		Address: "elys1ks3capc7dccnckn0crxgd2fppwgx0fup7j9phe",
		Amount:  math.NewInt(3300000),
		Claimed: false,
	},
	{
		Address: "elys1lx34u6l4f4d03l69fse2jqrt2qlk05l2hltaky",
		Amount:  math.NewInt(3300000),
		Claimed: false,
	},
	{
		Address: "elys1lyn04xrnlx98hd9metqclw4ldmmfjzf3ls5hnm",
		Amount:  math.NewInt(3300000),
		Claimed: false,
	},
	{
		Address: "elys1t9g9s2p93t9r4nlw37zn88eg92gsy70hy7t9q6",
		Amount:  math.NewInt(3300000),
		Claimed: false,
	},
	{
		Address: "elys1vvf5dkjyeays63wxrlf8evfgq39xh2jgfwul37",
		Amount:  math.NewInt(3300000),
		Claimed: false,
	},
	{
		Address: "elys12u28v8q28ey4s5tl97jxme9hegg6092d0ft6pv",
		Amount:  math.NewInt(3200000),
		Claimed: false,
	},
	{
		Address: "elys15ddjjjstha4za293jx0d3aw5sl40gav5vcn8z2",
		Amount:  math.NewInt(3200000),
		Claimed: false,
	},
	{
		Address: "elys17hjccx5drxp9stawf0hk7wu6fdjum5zvq3gnpt",
		Amount:  math.NewInt(3200000),
		Claimed: false,
	},
	{
		Address: "elys18gxxvw6slc0t0uva4leunvlgj0zra6axmauerc",
		Amount:  math.NewInt(3200000),
		Claimed: false,
	},
	{
		Address: "elys1av90mrfjuxsjmt9a2wf5wdpc4d8uecrfahla5t",
		Amount:  math.NewInt(3200000),
		Claimed: false,
	},
	{
		Address: "elys1rlun9gggkdrkkh84ju6cwe6ktlxeeeqsyrl2yt",
		Amount:  math.NewInt(3200000),
		Claimed: false,
	},
	{
		Address: "elys1s4zca4vew96r9uhlp85x5k9u064v83mca9vcaq",
		Amount:  math.NewInt(3200000),
		Claimed: false,
	},
	{
		Address: "elys10mccpak5yqax4h6czjq7uxzj9kxcvrtuc5qsp3",
		Amount:  math.NewInt(3100000),
		Claimed: false,
	},
	{
		Address: "elys1dtrqrljzwga9rrqcpru5gfpa5fnxnfp3lp9z0f",
		Amount:  math.NewInt(3100000),
		Claimed: false,
	},
	{
		Address: "elys1e2a05sthem6edd6d9thgzja55elsqscr6zsjuz",
		Amount:  math.NewInt(3100000),
		Claimed: false,
	},
	{
		Address: "elys1prgk39vg93uest9nh7c72d0nsy2n0p9kf9hfqp",
		Amount:  math.NewInt(3100000),
		Claimed: false,
	},
	{
		Address: "elys10383u7qjtwgtnnnptkg4lr7d9g3avatjedtad9",
		Amount:  math.NewInt(3000000),
		Claimed: false,
	},
	{
		Address: "elys1feykwd54hzd8flrpk3hxnm5u946jpp3qu4mcyq",
		Amount:  math.NewInt(3000000),
		Claimed: false,
	},
	{
		Address: "elys1j4732q4gmd8dsq9hjvmnvf0zny787w9395uce6",
		Amount:  math.NewInt(3000000),
		Claimed: false,
	},
	{
		Address: "elys14s7u46umfzsfg3d6lhxe4tfymlp5v9qxklanrw",
		Amount:  math.NewInt(2900000),
		Claimed: false,
	},
	{
		Address: "elys165luefet5lm9udwxzvfzp9l2ca8jkvdms5h974",
		Amount:  math.NewInt(2900000),
		Claimed: false,
	},
	{
		Address: "elys1tsr3dlxwl50jezu02ygccqa3jf3ctj6ht7j8d7",
		Amount:  math.NewInt(2900000),
		Claimed: false,
	},
	{
		Address: "elys1vd2dmqrzvhgk49hkw5kxv9zqlw3663ygvq6x2n",
		Amount:  math.NewInt(2900000),
		Claimed: false,
	},
	{
		Address: "elys1wmcr8clfzjau39xn6c48yj6k33906ry5rzyhq0",
		Amount:  math.NewInt(2900000),
		Claimed: false,
	},
	{
		Address: "elys1rkykyclp0nczq0dmzwky9ulc5fxyf7f75f456d",
		Amount:  math.NewInt(2800000),
		Claimed: false,
	},
	{
		Address: "elys1vujz0jwzwd9ht6jjg6sq6qryksl54vrufvlmp5",
		Amount:  math.NewInt(2800000),
		Claimed: false,
	},
	{
		Address: "elys1z9vuufe55ty8zfxgz35u5hz2j8st63c6a8vgup",
		Amount:  math.NewInt(2800000),
		Claimed: false,
	},
	{
		Address: "elys16v09tzzchlr35929ugsc42u9zcdz6wzk5l2dqj",
		Amount:  math.NewInt(2700000),
		Claimed: false,
	},
	{
		Address: "elys18spq95vhh4wm7wyedwdkykq2j4jw2g4h5lya27",
		Amount:  math.NewInt(2700000),
		Claimed: false,
	},
	{
		Address: "elys1d5fwqpf686tp7h7nr9js83a0cf2nkudew6hjl6",
		Amount:  math.NewInt(2700000),
		Claimed: false,
	},
	{
		Address: "elys1dalqypymz0c0pmcyrkzcn3s873qhjrzhhg2hze",
		Amount:  math.NewInt(2700000),
		Claimed: false,
	},
	{
		Address: "elys1m6us4ylf3emtn9m7namjl40g3mylur9nd9clxw",
		Amount:  math.NewInt(2700000),
		Claimed: false,
	},
	{
		Address: "elys1x9pnc0lrv6fpg04qz8pl78dcd3u9sec6edf6u7",
		Amount:  math.NewInt(2700000),
		Claimed: false,
	},
	{
		Address: "elys133cnaftcksvfene4djnuephkx2rcf94keu88y9",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys133klnqaayth6qj58yy2hnkpkgr75ukl4362w6h",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys154w89ak3zaztdafwl8mjsk6jclz9xd0mn6rtwv",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys1552eyhrpl0jhyvxv22pdej7jn4s8cu04lp9emh",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys15tvaqx6ja8rwcqnc6472d8v265dhxhj3m6g2gy",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys18x47wpd3rfpsc6dta3j7mawutudg8e5ag2t7su",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys1fj8jjwhq4lmq66mw07u4knqvdarn8ss6dvyxzz",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys1fl7facvy6nmmltyk6e7ddmdu0ng5n0hfml6hed",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys1m8rg06xecsfkdq677wnhrc6s707tywupps2mxd",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys1mk8dnkwxfgpuam7jy7xra59r48zgc8wrc4nzer",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys1qwhw43ctsk2djs9l2repj9kxuntktyufq653ky",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys1r6uqe8vuhq3hd883h622240567xhd7637gt538",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys1xvlhsaafn9ll8vsfahw0x2hn44jx5lplancgu8",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys1z4rre6rhj0nerv3kypf3l22rk27k5ut3y57zsc",
		Amount:  math.NewInt(2600000),
		Claimed: false,
	},
	{
		Address: "elys19thcps5j30y7vtr8495cfqg6y6lsezvepe0796",
		Amount:  math.NewInt(2500000),
		Claimed: false,
	},
	{
		Address: "elys1pemt63zctpxtmarevrjh0mwrv9mapnpyuwxak3",
		Amount:  math.NewInt(2500000),
		Claimed: false,
	},
	{
		Address: "elys1ulnhl75j7srh6vnk253xkyk4c7zzmmx7cv3ulm",
		Amount:  math.NewInt(2500000),
		Claimed: false,
	},
	{
		Address: "elys1z05gcfyvn68axc04jzlq5cv2ugcj2snfavsucg",
		Amount:  math.NewInt(2500000),
		Claimed: false,
	},
	{
		Address: "elys10tgk57n3hw3vxngv79y9m0lzszynfzkx90nrfq",
		Amount:  math.NewInt(2400000),
		Claimed: false,
	},
	{
		Address: "elys1458psg52pug56s4pjspjqe0h48wgeplxn53lcs",
		Amount:  math.NewInt(2400000),
		Claimed: false,
	},
	{
		Address: "elys15ax689auvqqe4yg93k5e9w0wcra64jsck8ewt3",
		Amount:  math.NewInt(2400000),
		Claimed: false,
	},
	{
		Address: "elys1cx5yv7340c459kz2fuugaj3p0086jzvd45lpwg",
		Amount:  math.NewInt(2400000),
		Claimed: false,
	},
	{
		Address: "elys1h9w60pls9mftphq76yxmemwxc2a25ed5wr38lm",
		Amount:  math.NewInt(2400000),
		Claimed: false,
	},
	{
		Address: "elys1kl40ghpaq7p23q2accjvy3r5z3ut7gtucf7dy7",
		Amount:  math.NewInt(2400000),
		Claimed: false,
	},
	{
		Address: "elys1ufu68c7ypqmmyg60drjjdpsfdezxuqh4zfagcf",
		Amount:  math.NewInt(2400000),
		Claimed: false,
	},
	{
		Address: "elys1x0ap56nyxthupdfwj8gez668jgcfhny70pxt9m",
		Amount:  math.NewInt(2400000),
		Claimed: false,
	},
	{
		Address: "elys1xurg4sqkymkdy4nygur7gkgmjr5kfkquxlyv45",
		Amount:  math.NewInt(2400000),
		Claimed: false,
	},
	{
		Address: "elys1yhwgwj3re4mnm9ugz2vlvtvxnjgd6cq3r86mym",
		Amount:  math.NewInt(2400000),
		Claimed: false,
	},
	{
		Address: "elys1zpllxnv2pkyf4fkx5g26nldq38yugc8mqvm04l",
		Amount:  math.NewInt(2400000),
		Claimed: false,
	},
	{
		Address: "elys10lgt4n9n2clfv9eutr5vrg3rj50526matg97zm",
		Amount:  math.NewInt(2300000),
		Claimed: false,
	},
	{
		Address: "elys193z8g5zd38dqqyjdx7pkzvcgl0g5eet2lc9vre",
		Amount:  math.NewInt(2300000),
		Claimed: false,
	},
	{
		Address: "elys1r46a5sx66yeenvy9areu5hqty8ftp2tu3vwrgn",
		Amount:  math.NewInt(2300000),
		Claimed: false,
	},
	{
		Address: "elys12hd9c07xdlx34w0uy3hy0sdkwxkurw0xvzplp9",
		Amount:  math.NewInt(2200000),
		Claimed: false,
	},
	{
		Address: "elys1custkwd5tr876ll3u99zx6tayuwh2ug8f02myf",
		Amount:  math.NewInt(2200000),
		Claimed: false,
	},
	{
		Address: "elys1kul5hljs77va07f92syyz948vh4gj7m2qq4k6g",
		Amount:  math.NewInt(2200000),
		Claimed: false,
	},
	{
		Address: "elys1q732rvh3h84r3duy5v53a9mymastmdvukdjrll",
		Amount:  math.NewInt(2200000),
		Claimed: false,
	},
	{
		Address: "elys1xnv34yz2f6c9zetktv6kd573856dvfck50l34m",
		Amount:  math.NewInt(2200000),
		Claimed: false,
	},
	{
		Address: "elys1ycr9dmcwkynaaqsh55pjupr2jzum9cmllgt2su",
		Amount:  math.NewInt(2200000),
		Claimed: false,
	},
	{
		Address: "elys16x5v5gkzx70xq0068erpyerymqf8gp4u0ksua9",
		Amount:  math.NewInt(2100000),
		Claimed: false,
	},
	{
		Address: "elys1cj7julquvmjj79tvmhat6ezn9tk7m4kvj46eqv",
		Amount:  math.NewInt(2100000),
		Claimed: false,
	},
	{
		Address: "elys1d85zeq8qekwz2smptnsph5qxfrls0398y9k7f5",
		Amount:  math.NewInt(2100000),
		Claimed: false,
	},
	{
		Address: "elys1d972cjcxyvz3u87lg4003p2gfv20902xrld792",
		Amount:  math.NewInt(2100000),
		Claimed: false,
	},
	{
		Address: "elys1lthtskrnh0kl93j7jldu0w2hq4etg5jrqvvk3x",
		Amount:  math.NewInt(2100000),
		Claimed: false,
	},
	{
		Address: "elys1rhy8zu4qtrkuk0gq6q4ddzeae2dd3lkqvjqvv7",
		Amount:  math.NewInt(2100000),
		Claimed: false,
	},
	{
		Address: "elys1rvdega423phx334fuupr3j6rend4swjazw395r",
		Amount:  math.NewInt(2100000),
		Claimed: false,
	},
	{
		Address: "elys1466n4d2044hng69q3f7el235clat3fdwapfgqa",
		Amount:  math.NewInt(2000000),
		Claimed: false,
	},
	{
		Address: "elys1gwmjefdskfx0ltlyrwxswyj7xlggtwr50zzxrn",
		Amount:  math.NewInt(2000000),
		Claimed: false,
	},
	{
		Address: "elys1l8urmk4gdssdmypd562tn8em0lgsvsyculzexv",
		Amount:  math.NewInt(2000000),
		Claimed: false,
	},
	{
		Address: "elys13apanp6v07zpyrgj9cp4t3p6zaq9fql4zzatlr",
		Amount:  math.NewInt(1900000),
		Claimed: false,
	},
	{
		Address: "elys15vkhfk7s6xhpm7dgwtyncal3x845ljzg2rduvv",
		Amount:  math.NewInt(1900000),
		Claimed: false,
	},
	{
		Address: "elys1eg6yd86unuutlt27p4cr9r8msju082q6gjyctc",
		Amount:  math.NewInt(1900000),
		Claimed: false,
	},
	{
		Address: "elys1lfa06ttvjxqeqewadxkkrt4qd98zf26ql3ap78",
		Amount:  math.NewInt(1900000),
		Claimed: false,
	},
	{
		Address: "elys1lsm7ughmnpma4ef969t58lnkkykyjj93zu9dgj",
		Amount:  math.NewInt(1900000),
		Claimed: false,
	},
	{
		Address: "elys1ny0ruas4trdgpjmldjsg73k9t9qftrl7hpw3zv",
		Amount:  math.NewInt(1900000),
		Claimed: false,
	},
	{
		Address: "elys1tugm892unf8w2hm0u0pc04xzrmplktgw3g6nl4",
		Amount:  math.NewInt(1900000),
		Claimed: false,
	},
	{
		Address: "elys1demwhtjenpw054hr7nz90wcnlsnjx77lln2rkd",
		Amount:  math.NewInt(1800000),
		Claimed: false,
	},
	{
		Address: "elys1dsyj0u9qqm3su7nmpsy0xlnsc77yzp8nsnlzgz",
		Amount:  math.NewInt(1800000),
		Claimed: false,
	},
	{
		Address: "elys148pkyu9q5ksuncvtn9rpv907ch8a4djqm64ge0",
		Amount:  math.NewInt(1700000),
		Claimed: false,
	},
	{
		Address: "elys155ln5e2udp52y8e2z8mhluv69z3grtkysf3e5t",
		Amount:  math.NewInt(1700000),
		Claimed: false,
	},
	{
		Address: "elys1sd8mj27rul8qd0c374zdd24h4p90wkyfx9fen4",
		Amount:  math.NewInt(1700000),
		Claimed: false,
	},
	{
		Address: "elys12uscvs73nceqhd56veegmtvqplegw0wy5snelm",
		Amount:  math.NewInt(1600000),
		Claimed: false,
	},
	{
		Address: "elys157fwtcj0l3kcn9se4sddm0z4unnkmhmyl3nn0k",
		Amount:  math.NewInt(1600000),
		Claimed: false,
	},
	{
		Address: "elys16x4ula2tu8cyw6amr4g2f9pt4txm0la7cruu3l",
		Amount:  math.NewInt(1600000),
		Claimed: false,
	},
	{
		Address: "elys17c6nnhet8eq7e3jqg3cxua444kgukz7gcqemej",
		Amount:  math.NewInt(1600000),
		Claimed: false,
	},
	{
		Address: "elys1k4akl8atf5j2sh86gqjd7h7z0vgm00p27zumx3",
		Amount:  math.NewInt(1600000),
		Claimed: false,
	},
	{
		Address: "elys1zscqmnt3gj0css4cc6033c73ewg6yeype7q2ev",
		Amount:  math.NewInt(1600000),
		Claimed: false,
	},
	{
		Address: "elys16gaqc4cpgsdjc0s9zp3pc08y4yvktt9qgfrcsa",
		Amount:  math.NewInt(1500000),
		Claimed: false,
	},
	{
		Address: "elys17m4q7xxg9khtuqm5slccu7pd8379uj335jj4n4",
		Amount:  math.NewInt(1500000),
		Claimed: false,
	},
	{
		Address: "elys17m9ues4fnpqmdgh8f49ms2hf97uw3vuzm6rhud",
		Amount:  math.NewInt(1500000),
		Claimed: false,
	},
	{
		Address: "elys19lffhg9xyy3gxs9s2lq23f62jaz4rpt0at756q",
		Amount:  math.NewInt(1500000),
		Claimed: false,
	},
	{
		Address: "elys1fxuzpj3ufnyzy8gstt4r3z5lfx0sx9lcvdn73e",
		Amount:  math.NewInt(1500000),
		Claimed: false,
	},
	{
		Address: "elys1g4mvn9f0anyara4myttygxwmk7hl8vk03353cw",
		Amount:  math.NewInt(1500000),
		Claimed: false,
	},
	{
		Address: "elys1gl3knus3m3u7jqj5ng2ymnc4exxghjqx4aft35",
		Amount:  math.NewInt(1500000),
		Claimed: false,
	},
	{
		Address: "elys10s3fhcnc0maqujef4pk0d8yduqqxphhz2vlpw7",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys14dmmzhm4fuvtdjpmat3ywyvpga7plq9d38f5mn",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys152mztn79tms5y3yyza9dds32vjjygr6d2z0xpc",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys186989e7j920kcpyew5w5ugsvdnc9kncxm0583g",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys1d4jrqxcfzl9eqk8mpqzeuml4ggsx2ap2ef6s7r",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys1eknx6ll55qvd3th7glwn0xzqssjk49ctsrmv44",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys1f34ymxphg7sas0nvxj5mj34zdwq4jrfcn45g22",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys1kqeqnh97t0we8fgaqxh3cgaueq5hvr3v9dsvk6",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys1ne2d8vmdpdfcgykrxm2jrdz893allefwv2l968",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys1nnqrg9zj2mtxgqeuqeht37utmnxgyx75pgvlpu",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys1rz9p7ccnpczakgmdkyrwajvjzaa2fsep6cpp3l",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys1u8hjvq2puecas2hgsrm8t5a5648kpq8n9hfnkk",
		Amount:  math.NewInt(1400000),
		Claimed: false,
	},
	{
		Address: "elys13nf72ahpfeewdgq8zratu094ynp50704n292gq",
		Amount:  math.NewInt(1300000),
		Claimed: false,
	},
	{
		Address: "elys14lpsl9qkqrarqh6hveawk7deqquvual7sqp00a",
		Amount:  math.NewInt(1300000),
		Claimed: false,
	},
	{
		Address: "elys1e0pv2a9gqqdvr2exs3x26llw7wj98t7wettfy9",
		Amount:  math.NewInt(1300000),
		Claimed: false,
	},
	{
		Address: "elys1f0urc2u0dd3ne6yl28armthxkfhz39szlklqxw",
		Amount:  math.NewInt(1300000),
		Claimed: false,
	},
	{
		Address: "elys1mztkd4c2xapd5crngml934fn5ajcd7duvka3r6",
		Amount:  math.NewInt(1300000),
		Claimed: false,
	},
	{
		Address: "elys1w8fn83xrz7c8l0h5rusu4f24h3h85q4kpmp9jv",
		Amount:  math.NewInt(1300000),
		Claimed: false,
	},
	{
		Address: "elys10zfug587fkgx99anww47v5xueu999dnskjewzj",
		Amount:  math.NewInt(1200000),
		Claimed: false,
	},
	{
		Address: "elys12y9fenuaqzusmyykkqrfjkya3zlkf0ac4nk3ls",
		Amount:  math.NewInt(1200000),
		Claimed: false,
	},
	{
		Address: "elys175599plkwxn2lgarsqjjxuk59z8pd07kddy2wu",
		Amount:  math.NewInt(1200000),
		Claimed: false,
	},
	{
		Address: "elys17hxh59zujdzqk3kqjn3mj9scsrs60yxagdwksz",
		Amount:  math.NewInt(1200000),
		Claimed: false,
	},
	{
		Address: "elys1gsp7kukpm8vrf4xng48emxpduf0q945wr0x2u0",
		Amount:  math.NewInt(1200000),
		Claimed: false,
	},
	{
		Address: "elys1lfd0ax27mx2s455gk7lzt078zkatp7jxr5h9vq",
		Amount:  math.NewInt(1200000),
		Claimed: false,
	},
	{
		Address: "elys12euywa7npkwje3xd9tcvc64ye42vuknaddxqx4",
		Amount:  math.NewInt(1100000),
		Claimed: false,
	},
	{
		Address: "elys19nh2u75czpm238n68kc9mvhcfsw0qyyqeufmu4",
		Amount:  math.NewInt(1100000),
		Claimed: false,
	},
	{
		Address: "elys1nlrl8s5mkd5va8yhv85gzk39gv4z7rp69chsdz",
		Amount:  math.NewInt(1100000),
		Claimed: false,
	},
	{
		Address: "elys1rkq90d9szx4wuwmjg39uuuwd7t4h7m65zessew",
		Amount:  math.NewInt(1100000),
		Claimed: false,
	},
	{
		Address: "elys1suhhaxnn2s5r64keypkdy43q6h2e075snc4aaj",
		Amount:  math.NewInt(1100000),
		Claimed: false,
	},
	{
		Address: "elys12rwjlvmgua5dscfxju98dexlkw989wz5ytr4kw",
		Amount:  math.NewInt(1000000),
		Claimed: false,
	},
	{
		Address: "elys19kny942rapnvhdhvpjuxfvl9w5dq83qkf2de5e",
		Amount:  math.NewInt(1000000),
		Claimed: false,
	},
	{
		Address: "elys1nahkrrrzfddr9rl6vxngua7rk4wt4w3l0qd4s6",
		Amount:  math.NewInt(1000000),
		Claimed: false,
	},
	{
		Address: "elys1v067755m9nyrf4e7wwzud7jk5hc9md94d5ats7",
		Amount:  math.NewInt(1000000),
		Claimed: false,
	},
	{
		Address: "elys1z99rmzum4hc3m2mtxk2gkrecvqthe3hyz07hgj",
		Amount:  math.NewInt(1000000),
		Claimed: false,
	},
	{
		Address: "elys12c7pvmhv7wehv7fr0m5u944keyftyn2kkyy69e",
		Amount:  math.NewInt(900000),
		Claimed: false,
	},
	{
		Address: "elys1h34nq0pvqwsd7a6z24mu9j2444s5gsln8gzukd",
		Amount:  math.NewInt(900000),
		Claimed: false,
	},
	{
		Address: "elys1jgn6jj76vmvuyudaxm0x35xtyp8nupa0pn96sc",
		Amount:  math.NewInt(900000),
		Claimed: false,
	},
	{
		Address: "elys1l7c7qx6vruu22c8av2am7cg40xx4rylpzcyrfv",
		Amount:  math.NewInt(900000),
		Claimed: false,
	},
	{
		Address: "elys1syngx88qhpxt0svyjgqyn6m7mcrqghu0t4dj9e",
		Amount:  math.NewInt(900000),
		Claimed: false,
	},
	{
		Address: "elys1x5ek7mlgem0y3nhhru5dcgnnpnk9qpdpezjj5j",
		Amount:  math.NewInt(900000),
		Claimed: false,
	},
	{
		Address: "elys1ypn06qf4aj5yjl5yv32xpa0d8zgwg59uyaccce",
		Amount:  math.NewInt(900000),
		Claimed: false,
	},
	{
		Address: "elys13gmuzyfqttajwxuqayefx5wp6zygv6nq7uwje7",
		Amount:  math.NewInt(800000),
		Claimed: false,
	},
	{
		Address: "elys14ze58see8spaq7v0gg6mk7rzl5kpcpaykmyjus",
		Amount:  math.NewInt(800000),
		Claimed: false,
	},
	{
		Address: "elys1c88rledwh3m3kkq76t66dyfaxz2hmxps23euhj",
		Amount:  math.NewInt(800000),
		Claimed: false,
	},
	{
		Address: "elys1ew9kncr6njm4qu489cq9qjwqnpczqrqy7ss3kx",
		Amount:  math.NewInt(800000),
		Claimed: false,
	},
	{
		Address: "elys1jwr238u9peh72u5t7xf6k2td28fetxx7zt3c58",
		Amount:  math.NewInt(800000),
		Claimed: false,
	},
	{
		Address: "elys1ldneed9628hs3qlap50hrt9h4skgd6s6mj29hz",
		Amount:  math.NewInt(800000),
		Claimed: false,
	},
	{
		Address: "elys1rh0u0gacxn9umc3qf4gwrsp40xfmz4veg9pkrf",
		Amount:  math.NewInt(800000),
		Claimed: false,
	},
	{
		Address: "elys1xy3ta4m3cvp9x9qzdan6nlt2e0h0tn0pawa0nh",
		Amount:  math.NewInt(800000),
		Claimed: false,
	},
	{
		Address: "elys15xxhnzdktr2at2tpl2mgth52sv4lv6u4d3k4nc",
		Amount:  math.NewInt(700000),
		Claimed: false,
	},
	{
		Address: "elys1fgkghg99mm7p97vyhsgq206qv703ujjemla2pa",
		Amount:  math.NewInt(700000),
		Claimed: false,
	},
	{
		Address: "elys1w74h5kh4q4h664hvrjt0rhttxfskrgk8xccqgc",
		Amount:  math.NewInt(700000),
		Claimed: false,
	},
	{
		Address: "elys17qf3avrea5dpf306pal9lt2k8vegqjnurfmm9q",
		Amount:  math.NewInt(600000),
		Claimed: false,
	},
	{
		Address: "elys1ae6sq4wdregwlx9k86pfk7aw3ck62z77fy6q8t",
		Amount:  math.NewInt(600000),
		Claimed: false,
	},
	{
		Address: "elys1c35ad98af8tnygtcplgak22dqpupl0ru0m4ffh",
		Amount:  math.NewInt(600000),
		Claimed: false,
	},
	{
		Address: "elys1d5wn9jrzjnzd6tc62qgrs7f452wxymhlht2ly7",
		Amount:  math.NewInt(600000),
		Claimed: false,
	},
	{
		Address: "elys1lraqygs45xzkme7zy8wrvmt67zdyxh7v8nrx2w",
		Amount:  math.NewInt(600000),
		Claimed: false,
	},
	{
		Address: "elys1mam3s73qyqedrjsfhzcs3c0hldyhsu02r0g8d2",
		Amount:  math.NewInt(600000),
		Claimed: false,
	},
	{
		Address: "elys1sq9xpwhzpw03glewq0trznjpz93tshauvme503",
		Amount:  math.NewInt(600000),
		Claimed: false,
	},
	{
		Address: "elys13dnvp7dl9zlpdt99pn3kqckk3va6thysq28zz8",
		Amount:  math.NewInt(500000),
		Claimed: false,
	},
	{
		Address: "elys14ux0shuj3hyqc2ywvn0hjud5q0lrswn9qrqr74",
		Amount:  math.NewInt(500000),
		Claimed: false,
	},
	{
		Address: "elys17mzu94uf6k0c5hg7zjdw0x333ts6x4vlufue8g",
		Amount:  math.NewInt(500000),
		Claimed: false,
	},
	{
		Address: "elys1jrhhhkypmppgcj73ptzm0zn3dnsj2823sh9tm5",
		Amount:  math.NewInt(500000),
		Claimed: false,
	},
	{
		Address: "elys1mkppcy9vuyjkj0n8jyauf297tur9rgt392jftw",
		Amount:  math.NewInt(500000),
		Claimed: false,
	},
	{
		Address: "elys1qr6z0svzhz0gzg7ja0chly5t2n64nyjrytfnzq",
		Amount:  math.NewInt(500000),
		Claimed: false,
	},
	{
		Address: "elys1t5fz4wgqdzssvmax9r40wd2uwgjwrtkkmpxfc2",
		Amount:  math.NewInt(500000),
		Claimed: false,
	},
	{
		Address: "elys1xr9u3tynw8vwkwccdq8af27rnd9ggvpk78q696",
		Amount:  math.NewInt(500000),
		Claimed: false,
	},
	{
		Address: "elys1y250mef7xe5ekskhw4pdjdemkm3czakmyqpnt0",
		Amount:  math.NewInt(500000),
		Claimed: false,
	},
	{
		Address: "elys12pjgwmk672r98fhjrpxg9xz2j29vmu9rm70xw2",
		Amount:  math.NewInt(400000),
		Claimed: false,
	},
	{
		Address: "elys14hmzeg0jsm7sty59vxqvcpm7wpuc3q0pz490x0",
		Amount:  math.NewInt(400000),
		Claimed: false,
	},
	{
		Address: "elys19dl6gws9t0w2qrme9k2n48v8mpargdrsf0228r",
		Amount:  math.NewInt(400000),
		Claimed: false,
	},
	{
		Address: "elys1ecgc6jm3pyk3wugczkx4gc4cq5r467fxz280t7",
		Amount:  math.NewInt(400000),
		Claimed: false,
	},
	{
		Address: "elys1hskzcjss4l7hqdsfdq2wn7j9w8trmun8x38rjm",
		Amount:  math.NewInt(400000),
		Claimed: false,
	},
	{
		Address: "elys1j4cc5f7mvgw3lhyeswdwk3pqtcj5l85q99950r",
		Amount:  math.NewInt(400000),
		Claimed: false,
	},
	{
		Address: "elys1unuzyq0d2r7zn5r8m222cjqygg85ggmdff6vk2",
		Amount:  math.NewInt(400000),
		Claimed: false,
	},
	{
		Address: "elys164rshrwr7h035m4hy7yvvlg2gav8yn0az3zf58",
		Amount:  math.NewInt(300000),
		Claimed: false,
	},
	{
		Address: "elys182vtppkgzqkgyu7gf5aur8jgrxq0g5l87cfaq8",
		Amount:  math.NewInt(300000),
		Claimed: false,
	},
	{
		Address: "elys1tdtxrmgyelzswvd88w9uec8h307hpwahr0xgvh",
		Amount:  math.NewInt(300000),
		Claimed: false,
	},
	{
		Address: "elys1z5snzwwmel7u44p6awh29e9sh6agvvx7dcq922",
		Amount:  math.NewInt(300000),
		Claimed: false,
	},
	{
		Address: "elys10xm5qaqcsfatgksreepg8nu09ltkmjkdu8t79g",
		Amount:  math.NewInt(200000),
		Claimed: false,
	},
	{
		Address: "elys15gwq0mrw67jangvund48gkfra4eq0y2ucquu3p",
		Amount:  math.NewInt(200000),
		Claimed: false,
	},
	{
		Address: "elys1xw24v52sw468z8qan7vt9vhzt9qegu4y4m3cqv",
		Amount:  math.NewInt(200000),
		Claimed: false,
	},
	{
		Address: "elys16l9f33rfpqv6rknkj6dza00804j6rfkedte79r",
		Amount:  math.NewInt(100000),
		Claimed: false,
	},
	{
		Address: "elys17xqffuqq24000egypte2l6pyn56jwm7tt6ftks",
		Amount:  math.NewInt(100000),
		Claimed: false,
	},
	{
		Address: "elys1f6t08wjat6qvwl9ek25uh42p0csn92r2qq3pah",
		Amount:  math.NewInt(100000),
		Claimed: false,
	},
	{
		Address: "elys1gn7zgx66e7k2fjex45kdwnpn73jl3nyh0wnwa3",
		Amount:  math.NewInt(100000),
		Claimed: false,
	},
	{
		Address: "elys1jcw9lta8fseqn2he9za2mgyvefjgflpfumdjde",
		Amount:  math.NewInt(100000),
		Claimed: false,
	},
	{
		Address: "elys1l2r20fh98shncew80cv3rne83kf4997g7fumxn",
		Amount:  math.NewInt(100000),
		Claimed: false,
	},
	{
		Address: "elys1qej7dvr5spd8gt698a5sesa0seuwq4dmr62vs9",
		Amount:  math.NewInt(100000),
		Claimed: false,
	},
	{
		Address: "elys1w00em8a5uk6x8neh746v2p27kk9vn8c850uk02",
		Amount:  math.NewInt(100000),
		Claimed: false,
	},
}
